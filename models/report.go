package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/mbvlabs/plyo-hackathon/models/internal/db"
)

type Report struct {
	ID                               uuid.UUID
	CreatedAt                        time.Time
	UpdatedAt                        time.Time
	CompanyCandidateID               string
	CompanyName                      string
	Status                           string
	ProgressPercentage               int64
	PreliminaryResearchCompleted     bool
	CompanyIntelligenceCompleted     bool
	CompetitiveIntelligenceCompleted bool
	MarketDynamicsCompleted          bool
	TrendAnalysisCompleted           bool
	CompanyIntelligenceData          string
	CompetitiveIntelligenceData      string
	MarketDynamicsData               string
	TrendAnalysisData                string
	FinalReport                      string
	CompletedAt                      time.Time
}

func FindReport(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Report, error) {
	row, err := db.New().QueryReportByID(ctx, dbtx, id.String())
	if err != nil {
		return Report{}, err
	}

	result, err := rowToReport(row)
	if err != nil {
		return Report{}, err
	}
	return result, nil
}

type CreateReportData struct {
	CompanyCandidateID               string
	CompanyName                      string
	Status                           string
	ProgressPercentage               int64
	PreliminaryResearchCompleted     bool
	CompanyIntelligenceCompleted     bool
	CompetitiveIntelligenceCompleted bool
	MarketDynamicsCompleted          bool
	TrendAnalysisCompleted           bool
	CompanyIntelligenceData          string
	CompetitiveIntelligenceData      string
	MarketDynamicsData               string
	TrendAnalysisData                string
	FinalReport                      string
	CompletedAt                      time.Time
}

func CreateReport(
	ctx context.Context,
	dbtx db.DBTX,
	data CreateReportData,
) (Report, error) {
	if err := validate.Struct(data); err != nil {
		return Report{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.NewInsertReportParams(
		data.CompanyCandidateID,
		data.CompanyName,
		data.Status,
		sql.NullInt64{Int64: data.ProgressPercentage, Valid: true},
		sql.NullBool{Bool: data.PreliminaryResearchCompleted, Valid: true},
		sql.NullBool{Bool: data.CompanyIntelligenceCompleted, Valid: true},
		sql.NullBool{Bool: data.CompetitiveIntelligenceCompleted, Valid: true},
		sql.NullBool{Bool: data.MarketDynamicsCompleted, Valid: true},
		sql.NullBool{Bool: data.TrendAnalysisCompleted, Valid: true},
		sql.NullString{String: data.CompanyIntelligenceData, Valid: true},
		sql.NullString{String: data.CompetitiveIntelligenceData, Valid: true},
		sql.NullString{String: data.MarketDynamicsData, Valid: true},
		sql.NullString{String: data.TrendAnalysisData, Valid: true},
		sql.NullString{String: data.FinalReport, Valid: true},
		sql.NullTime{Time: data.CompletedAt, Valid: true},
	)
	row, err := db.New().InsertReport(ctx, dbtx, params)
	if err != nil {
		return Report{}, err
	}

	result, err := rowToReport(row)
	if err != nil {
		return Report{}, err
	}
	return result, nil
}

func DestroyReport(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.New().DeleteReport(ctx, dbtx, id.String())
}

func AllReports(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Report, error) {
	rows, err := db.New().QueryAllReports(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	reports := make([]Report, len(rows))
	for i, row := range rows {
		result, err := rowToReport(row)
		if err != nil {
			return nil, err
		}
		reports[i] = result
	}

	return reports, nil
}

type PaginatedReports struct {
	Reports    []Report
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func PaginateReports(
	ctx context.Context,
	dbtx db.DBTX,
	page int64,
	pageSize int64,
) (PaginatedReports, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := db.New().CountReports(ctx, dbtx)
	if err != nil {
		return PaginatedReports{}, err
	}

	rows, err := db.New().QueryPaginatedReports(
		ctx,
		dbtx,
		db.NewQueryPaginatedReportsParams(pageSize, offset),
	)
	if err != nil {
		return PaginatedReports{}, err
	}

	reports := make([]Report, len(rows))
	for i, row := range rows {
		result, err := rowToReport(row)
		if err != nil {
			return PaginatedReports{}, err
		}
		reports[i] = result
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedReports{
		Reports:    reports,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func UpdateCompanyIntelligence(
	ctx context.Context,
	dbtx db.DBTX,
	reportID uuid.UUID,
	data string,
) error {
	return db.New().UpdateCompanyIntelligence(ctx, dbtx, db.UpdateCompanyIntelligenceParams{
		CompanyIntelligenceData:      sql.NullString{Valid: true, String: data},
		CompanyIntelligenceCompleted: sql.NullBool{Bool: true, Valid: true},
		ID:                           reportID.String(),
	})
}

func UpdateCompetitiveIntelligence(
	ctx context.Context,
	dbtx db.DBTX,
	reportID uuid.UUID,
	data string,
) error {
	return db.New().
		UpdateCompetitiveIntelligence(ctx, dbtx, db.UpdateCompetitiveIntelligenceParams{
			CompetitiveIntelligenceData:      sql.NullString{String: data, Valid: true},
			CompetitiveIntelligenceCompleted: sql.NullBool{Bool: true, Valid: true},
			ID:                               reportID.String(),
		})
}

func UpdateMarketDynamics(
	ctx context.Context,
	dbtx db.DBTX,
	reportID uuid.UUID,
	data string,
) error {
	return db.New().UpdateMarketDynamics(ctx, dbtx, db.UpdateMarketDynamicsParams{
		MarketDynamicsData:      sql.NullString{String: data, Valid: true},
		MarketDynamicsCompleted: sql.NullBool{Bool: true, Valid: true},
		ID:                      reportID.String(),
	})
}

func UpdateTrendAnalysis(
	ctx context.Context,
	dbtx db.DBTX,
	reportID uuid.UUID,
	data string,
) error {
	return db.New().UpdateTrendAnalysis(ctx, dbtx, db.UpdateTrendAnalysisParams{
		TrendAnalysisData:      sql.NullString{String: data, Valid: true},
		TrendAnalysisCompleted: sql.NullBool{Bool: true, Valid: true},
		ID:                     reportID.String(),
	})
}

func CalculateProgress(report Report) int64 {
	completedCount := int64(0)
	totalAgents := int64(4)

	if report.CompanyIntelligenceCompleted {
		completedCount++
	}
	if report.CompetitiveIntelligenceCompleted {
		completedCount++
	}
	if report.MarketDynamicsCompleted {
		completedCount++
	}
	if report.TrendAnalysisCompleted {
		completedCount++
	}

	progress := (completedCount * 100) / totalAgents

	// If all agents are completed and final report exists, show 100%
	// This prevents the progress from dropping when final report is being generated
	if completedCount == totalAgents && report.FinalReport != "" {
		return 100
	}

	return progress
}

func UpdateReportProgressToStarted(
	ctx context.Context,
	dbtx db.DBTX,
	reportID uuid.UUID,
) error {
	status := "in_progress"

	return db.New().UpdateReportProgress(ctx, dbtx, db.UpdateReportProgressParams{
		ProgressPercentage: sql.NullInt64{Int64: 0, Valid: true},
		Status:             status,
		ID:                 reportID.String(),
	})
}

func UpdateReportProgress(
	ctx context.Context,
	dbtx db.DBTX,
	reportID uuid.UUID,
) error {
	report, err := FindReport(ctx, dbtx, reportID)
	if err != nil {
		return err
	}

	progress := CalculateProgress(report)
	status := "in_progress"

	switch {
	case progress == 100:
		status = "completed"
	case progress > 0 && progress < 100:
		status = "processing"
	case progress == 0:
		status = "pending"
	}

	return db.New().UpdateReportProgress(ctx, dbtx, db.UpdateReportProgressParams{
		ProgressPercentage: sql.NullInt64{Int64: progress, Valid: true},
		Status:             status,
		ID:                 reportID.String(),
	})
}

func UpdateFinalReport(
	ctx context.Context,
	dbtx db.DBTX,
	reportID uuid.UUID,
	finalReport string,
) error {
	return db.New().UpdateFinalReport(ctx, dbtx, db.UpdateFinalReportParams{
		FinalReport: sql.NullString{String: finalReport, Valid: true},
		ID:          reportID.String(),
	})
}

func rowToReport(row db.Report) (Report, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return Report{}, err
	}

	return Report{
		ID:                               id,
		CreatedAt:                        row.CreatedAt,
		UpdatedAt:                        row.UpdatedAt,
		CompanyCandidateID:               row.CompayCandidateID,
		CompanyName:                      row.CompanyName,
		Status:                           row.Status,
		ProgressPercentage:               row.ProgressPercentage.Int64,
		PreliminaryResearchCompleted:     row.PreliminaryResearchCompleted.Bool,
		CompanyIntelligenceCompleted:     row.CompanyIntelligenceCompleted.Bool,
		CompetitiveIntelligenceCompleted: row.CompetitiveIntelligenceCompleted.Bool,
		MarketDynamicsCompleted:          row.MarketDynamicsCompleted.Bool,
		TrendAnalysisCompleted:           row.TrendAnalysisCompleted.Bool,
		CompanyIntelligenceData:          row.CompanyIntelligenceData.String,
		CompetitiveIntelligenceData:      row.CompetitiveIntelligenceData.String,
		MarketDynamicsData:               row.MarketDynamicsData.String,
		TrendAnalysisData:                row.TrendAnalysisData.String,
		FinalReport:                      row.FinalReport.String,
		CompletedAt:                      row.CompletedAt.Time,
	}, nil
}
