package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/mbvlabs/plyo-hackathon/agents"
	"github.com/mbvlabs/plyo-hackathon/database"
	"github.com/mbvlabs/plyo-hackathon/models"
	"github.com/mbvlabs/plyo-hackathon/router/cookies"
	"github.com/mbvlabs/plyo-hackathon/views"
	"maragu.dev/goqite"
	"maragu.dev/goqite/jobs"
)

type Reports struct {
	db database.SQLite
	q  *goqite.Queue
}

func newReports(
	db database.SQLite,
	q *goqite.Queue,
) Reports {
	return Reports{db, q}
}

func (r Reports) Create(c echo.Context) error {
	id := c.QueryParam("id")

	companyCandidateID, err := uuid.Parse(id)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"invalid company candidate ID",
			"error", err,
			"candidate_id", id,
		)
		return render(c, views.BadRequest())
	}

	candidate, err := models.FindCompanyCandidates(
		c.Request().Context(),
		r.db.Conn(),
		companyCandidateID,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to find company candidate",
			"error", err,
			"candidate_id", companyCandidateID,
		)
		return render(c, views.NotFound())
	}

	data := models.CreateReportData{
		CompanyCandidateID:               candidate.ID.String(),
		CompanyName:                      candidate.Name,
		Status:                           "pending",
		ProgressPercentage:               0,
		PreliminaryResearchCompleted:     true, // Since we already have the preliminary research
		CompanyIntelligenceCompleted:     false,
		CompetitiveIntelligenceCompleted: false,
		MarketDynamicsCompleted:          false,
		TrendAnalysisCompleted:           false,
		CompanyIntelligenceData:          "",
		CompetitiveIntelligenceData:      "",
		MarketDynamicsData:               "",
		TrendAnalysisData:                "",
		FinalReport:                      "",
		CompletedAt:                      time.Time{},
	}

	report, err := models.CreateReport(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to create report",
			"error", err,
		)
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to create report: %v", err)); flashErr != nil {
			return flashErr
		}
		return render(c, views.InternalError())
	}

	return getSSE(c).Redirect(fmt.Sprintf("/reports/%s", report.ID.String()))
}

func (r Reports) Show(c echo.Context) error {
	reportID := c.Param("id")

	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"invalid report ID",
			"error", err,
			"report_id", reportID,
		)
		return render(c, views.BadRequest())
	}

	report, err := models.FindReport(
		c.Request().Context(),
		r.db.Conn(),
		reportUUID,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to find report",
			"error", err,
			"report_id", reportUUID,
		)
		return render(c, views.NotFound())
	}

	if report.Status == "pending" {
		company, err := models.FindCompanyCandidates(
			c.Request().Context(),
			r.db.Conn(),
			uuid.MustParse(report.CompanyCandidateID),
		)
		if err != nil {
			return err
		}

		if err = models.UpdateReportProgress(
			c.Request().Context(),
			r.db.Conn(),
			reportUUID,
		); err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to update report progress",
				"error", err,
				"report_id", reportUUID,
			)
		}

		// INTEL
		params := agents.CompanyIntelligenceJobParams{
			ReportID:      report.ID,
			CandidateName: report.CompanyName,
			CompanyURL:    company.Domain,
		}
		data, err := json.Marshal(params)
		if err != nil {
			log.Info("Error marshalling job", "error", err)
		}

		if err := jobs.Create(c.Request().Context(), r.q, agents.CompanyIntelligenceJobName, data); err != nil {
			log.Info("Error creating job", "error", err)
		}

		// COMPETITORS
		competitiveParams := agents.CompetitiveIntelligenceJobParams{
			ReportID:      report.ID,
			CandidateName: report.CompanyName,
			CompanyURL:    company.Domain,
		}
		competitiveData, err := json.Marshal(competitiveParams)
		if err != nil {
			log.Info("Error marshalling job", "error", err)
		}

		if err := jobs.Create(c.Request().Context(), r.q, agents.CompetitiveIntelligenceJobName, competitiveData); err != nil {
			log.Info("Error creating job", "error", err)
		}

		// TREND
		trendParams := agents.TrendAnalysisJobParams{
			ReportID:      report.ID,
			CandidateName: report.CompanyName,
			CompanyURL:    company.Domain,
		}
		trendData, err := json.Marshal(trendParams)
		if err != nil {
			log.Info("Error marshalling job", "error", err)
		}

		if err := jobs.Create(c.Request().Context(), r.q, agents.TrendAnalysisJobName, trendData); err != nil {
			log.Info("Error creating job", "error", err)
		}

		// MARKET
		marketParams := agents.MarketDynamicsJobParams{
			ReportID:      report.ID,
			CandidateName: report.CompanyName,
			CompanyURL:    company.Domain,
		}
		marketData, err := json.Marshal(marketParams)
		if err != nil {
			log.Info("Error marshalling job", "error", err)
		}

		if err := jobs.Create(c.Request().Context(), r.q, agents.MarketDynamicsJobName, marketData); err != nil {
			log.Info("Error creating job", "error", err)
		}
	}

	slog.InfoContext(
		c.Request().Context(),
		"showing report chat",
		"report_id", report.ID,
		"company_name", report.CompanyName,
		"status", report.Status,
		"progress", report.ProgressPercentage,
	)

	return render(c, views.ReportChat(report))
}

func (r Reports) TrackReportProgress(c echo.Context) error {
	reportID := c.Param("id")

	// Parse report ID
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"invalid report ID",
			"error", err,
			"report_id", reportID,
		)
		return c.String(400, "Invalid report ID")
	}

	sse := getSSE(c)

	report, err := models.FindReport(
		c.Request().Context(),
		r.db.Conn(),
		reportUUID,
	)
	if err != nil {
		return c.String(404, "Report not found")
	}

	if allAgentsCompleted(report) && report.FinalReport == "" {
		company, err := models.FindCompanyCandidates(
			c.Request().Context(),
			r.db.Conn(),
			uuid.MustParse(report.CompanyCandidateID),
		)
		if err != nil {
			return err
		}

		reportGenParams := agents.ReportGeneratorJobParams{
			ReportID:      report.ID,
			CandidateName: report.CompanyName,
			CompanyURL:    company.Domain,
		}
		reportGenData, err := json.Marshal(reportGenParams)
		if err != nil {
			log.Info("Error marshalling job", "error", err)
		}

		if err := jobs.Create(c.Request().Context(), r.q, agents.ReportGeneratorJobName, reportGenData); err != nil {
			log.Info("Error creating job", "error", err)
		}

		sse.PatchElements("<div id='reportProgressPuller'></div>")
	}

	sse.PatchElementTempl(views.ReportProgress(report))
	sse.PatchElementTempl(views.ReportUpdated(report.UpdatedAt))
	return sse.PatchElementTempl(views.ReportHeaderProgress(report))
}

func allAgentsCompleted(report models.Report) bool {
	return report.CompanyIntelligenceCompleted &&
		report.CompetitiveIntelligenceCompleted &&
		report.MarketDynamicsCompleted &&
		report.TrendAnalysisCompleted
}
