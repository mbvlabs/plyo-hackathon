package controllers

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mbvlabs/plyo-hackathon/agents"
	"github.com/mbvlabs/plyo-hackathon/config"
	"github.com/mbvlabs/plyo-hackathon/database"
	"github.com/mbvlabs/plyo-hackathon/models"
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/router/cookies"
	"github.com/mbvlabs/plyo-hackathon/tools"
	"github.com/mbvlabs/plyo-hackathon/views"
	"github.com/starfederation/datastar-go/datastar"
)

type Reports struct {
	db database.SQLite
}

func newReports(db database.SQLite) Reports {
	return Reports{db}
}

// type CreateReportFormPayload struct {
// 	CompanyCandidateID string `query:"id"`
// }

func (r Reports) Create(c echo.Context) error {
	// var payload CreateReportFormPayload
	// if err := c.Bind(&payload); err != nil {
	// 	slog.ErrorContext(
	// 		c.Request().Context(),
	// 		"could not parse CreateReportFormPayload",
	// 		"error",
	// 		err,
	// 	)
	// 	return render(c, views.BadRequest())
	// }
	id := c.QueryParam("id")

	slog.Info("YOOOOOOOOOOOOOOOO", "id", id)

	// Parse company candidate ID
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

	// Find the company candidate
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

	slog.Info("CANDIDAAAAAAAAAAAAAAAAAATE", "candi", candidate)

	// Create the report
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

	slog.InfoContext(
		c.Request().Context(),
		"research report created",
		"report_id", report.ID,
		"company_name", candidate.Name,
	)

	// Redirect to the report chat page
	redirectURL := fmt.Sprintf("/reports/%s", report.ID.String())
	return getSSE(c).Redirect(redirectURL)
}

func (r Reports) Show(c echo.Context) error {
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
		return render(c, views.BadRequest())
	}

	// Find the report
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

	// Check if report is pending and should start processing
	if report.Status == "pending" {
		// Update status to in_progress to indicate processing has started
		err = models.UpdateReportProgress(
			c.Request().Context(),
			r.db.Conn(),
			reportUUID,
		)
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to update report progress",
				"error", err,
				"report_id", reportUUID,
			)
		} else {
			// Reload the report to get the updated status
			report, _ = models.FindReport(
				c.Request().Context(),
				r.db.Conn(),
				reportUUID,
			)
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

// StreamProgress handles Server-Sent Events for real-time progress updates
func (r Reports) StreamProgress(c echo.Context) error {
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

	// Send initial report state
	report, err := models.FindReport(
		c.Request().Context(),
		r.db.Conn(),
		reportUUID,
	)
	if err != nil {
		return c.String(404, "Report not found")
	}

	// Send initial state
	sse.PatchElementTempl(views.ReportProgress(report))

	// If report needs processing, start agents and stream progress
	if report.Status == "pending" ||
		(report.Status == "in_progress" && !allAgentsCompleted(report)) {
		return r.runAgentsWithStreaming(c, sse, report)
	}

	// If already completed, just return
	return nil
}

// Check if all agents are completed
func allAgentsCompleted(report models.Report) bool {
	return report.CompanyIntelligenceCompleted &&
		report.CompetitiveIntelligenceCompleted &&
		report.MarketDynamicsCompleted &&
		report.TrendAnalysisCompleted
}

// runAgentsWithStreaming runs the research agents while keeping SSE connection alive
func (r Reports) runAgentsWithStreaming(
	c echo.Context,
	sse *datastar.ServerSentEventGenerator,
	report models.Report,
) error {
	slog.Info("STAAAAAAAAAAARTING")
	ctx := c.Request().Context()

	// Get company candidate for URL
	candidateID, err := uuid.Parse(report.CompanyCandidateID)
	if err != nil {
		slog.ErrorContext(ctx, "invalid candidate ID", "error", err)
		return err
	}

	candidate, err := models.FindCompanyCandidates(ctx, r.db.Conn(), candidateID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find candidate", "error", err)
		return err
	}

	// Setup tools and provider
	serper := tools.NewSerper(config.App.SerperAPIkey)
	openai := providers.NewClient(config.App.OpenAPIKey)
	toolsMap := map[string]tools.Tooler{serper.GetName(): &serper}

	// Create agents
	companyIntel := agents.NewCompanyIntelligence(openai, toolsMap)
	competitiveIntel := agents.NewCompetitiveIntelligence(openai, toolsMap)
	marketDynamics := agents.NewMarketDynamics(openai, toolsMap)
	trendAnalysis := agents.NewTrendAnalysis(openai, toolsMap)

	// Build company URL
	companyURL := fmt.Sprintf("https://%s", candidate.Domain)

	// Create channels for agent completion
	companyDone := make(chan error, 1)
	competitiveDone := make(chan error, 1)
	marketDone := make(chan error, 1)
	trendDone := make(chan error, 1)

	// Start agents in parallel
	go func() {
		slog.InfoContext(ctx, "starting company intelligence", "report_id", report.ID)
		result, err := companyIntel.Research(ctx, candidate.Name, companyURL)
		if err != nil {
			slog.ErrorContext(ctx, "company intelligence failed", "error", err)
			companyDone <- err
			return
		}
		if err := models.UpdateCompanyIntelligence(ctx, r.db.Conn(), report.ID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update company intelligence", "error", err)
			companyDone <- err
			return
		}
		models.UpdateReportProgress(ctx, r.db.Conn(), report.ID)
		slog.InfoContext(ctx, "completed company intelligence", "report_id", report.ID)
		companyDone <- nil
	}()

	go func() {
		slog.InfoContext(ctx, "starting competitive intelligence", "report_id", report.ID)
		result, err := competitiveIntel.Research(ctx, candidate.Name, companyURL)
		if err != nil {
			slog.ErrorContext(ctx, "competitive intelligence failed", "error", err)
			competitiveDone <- err
			return
		}
		if err := models.UpdateCompetitiveIntelligence(ctx, r.db.Conn(), report.ID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update competitive intelligence", "error", err)
			competitiveDone <- err
			return
		}
		models.UpdateReportProgress(ctx, r.db.Conn(), report.ID)
		slog.InfoContext(ctx, "completed competitive intelligence", "report_id", report.ID)
		competitiveDone <- nil
	}()

	go func() {
		slog.InfoContext(ctx, "starting market dynamics", "report_id", report.ID)
		result, err := marketDynamics.Research(ctx, candidate.Name, companyURL)
		if err != nil {
			slog.ErrorContext(ctx, "market dynamics failed", "error", err)
			marketDone <- err
			return
		}
		if err := models.UpdateMarketDynamics(ctx, r.db.Conn(), report.ID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update market dynamics", "error", err)
			marketDone <- err
			return
		}
		models.UpdateReportProgress(ctx, r.db.Conn(), report.ID)
		slog.InfoContext(ctx, "completed market dynamics", "report_id", report.ID)
		marketDone <- nil
	}()

	go func() {
		slog.InfoContext(ctx, "starting trend analysis", "report_id", report.ID)
		result, err := trendAnalysis.Research(ctx, candidate.Name, companyURL)
		if err != nil {
			slog.ErrorContext(ctx, "trend analysis failed", "error", err)
			trendDone <- err
			return
		}
		if err := models.UpdateTrendAnalysis(ctx, r.db.Conn(), report.ID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update trend analysis", "error", err)
			trendDone <- err
			return
		}
		models.UpdateReportProgress(ctx, r.db.Conn(), report.ID)
		slog.InfoContext(ctx, "completed trend analysis", "report_id", report.ID)
		trendDone <- nil
	}()

	// Stream progress updates
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	completedAgents := 0
	totalAgents := 4

	for completedAgents < totalAgents {
		select {
		case <-ctx.Done():
			return nil
		case err := <-companyDone:
			completedAgents++
			if err != nil {
				slog.ErrorContext(ctx, "company intelligence agent failed", "error", err)
			}
		case err := <-competitiveDone:
			completedAgents++
			if err != nil {
				slog.ErrorContext(ctx, "competitive intelligence agent failed", "error", err)
			}
		case err := <-marketDone:
			completedAgents++
			if err != nil {
				slog.ErrorContext(ctx, "market dynamics agent failed", "error", err)
			}
		case err := <-trendDone:
			completedAgents++
			if err != nil {
				slog.ErrorContext(ctx, "trend analysis agent failed", "error", err)
			}
		case <-ticker.C:
			// Send periodic updates
			updatedReport, err := models.FindReport(ctx, r.db.Conn(), report.ID)
			if err == nil {
				sse.PatchElementTempl(views.ReportProgress(updatedReport))
				report = updatedReport
			}
		}

		// Send update after each agent completion
		if completedAgents > 0 {
			updatedReport, err := models.FindReport(ctx, r.db.Conn(), report.ID)
			if err == nil {
				sse.PatchElementTempl(views.ReportProgress(updatedReport))
				report = updatedReport
			}
		}
	}

	slog.InfoContext(ctx, "all research agents completed", "report_id", report.ID)

	// Generate final report using all collected data
	finalReport, err := models.FindReport(ctx, r.db.Conn(), report.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find report for final generation", "error", err)
		return err
	}

	// Create report generator and generate final report
	reportGenerator := agents.NewReportGenerator(openai, toolsMap)
	finalReportContent, err := reportGenerator.Generate(
		ctx,
		candidate.Name,
		companyURL,
		finalReport.CompanyIntelligenceData,
		finalReport.CompetitiveIntelligenceData,
		finalReport.MarketDynamicsData,
		finalReport.TrendAnalysisData,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate final report", "error", err)
		return err
	}

	// Update report with final content and mark as completed
	err = models.UpdateFinalReport(ctx, r.db.Conn(), report.ID, finalReportContent)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update final report", "error", err)
		return err
	}

	slog.InfoContext(ctx, "final report generated successfully", "report_id", report.ID)

	// Send final update
	completedReport, err := models.FindReport(ctx, r.db.Conn(), report.ID)
	if err == nil {
		sse.PatchElementTempl(views.ReportProgress(completedReport))
	}

	return nil
}
