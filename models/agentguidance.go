package models

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/mbvlabs/plyo-hackathon/models/internal/db"
)

type AgentGuidance struct {
	ID              uuid.UUID
	ResearchBriefID string
	GuidanceKey     string
	GuidanceValue   string
}

func FindAgentGuidance(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (AgentGuidance, error) {
	row, err := db.New().QueryAgentGuidanceByID(ctx, dbtx, id.String())
	if err != nil {
		return AgentGuidance{}, err
	}

	result, err := rowToAgentGuidance(row)
	if err != nil {
		return AgentGuidance{}, err
	}
	return result, nil
}

func FindAgentGuidancesByResearchBriefID(
	ctx context.Context,
	dbtx db.DBTX,
	researchBriefID string,
) ([]AgentGuidance, error) {
	rows, err := db.New().QueryAgentGuidancesByResearchBriefID(ctx, dbtx, researchBriefID)
	if err != nil {
		return nil, err
	}

	guidances := make([]AgentGuidance, len(rows))
	for i, row := range rows {
		result, err := rowToAgentGuidance(row)
		if err != nil {
			return nil, err
		}
		guidances[i] = result
	}

	return guidances, nil
}

type CreateAgentGuidanceData struct {
	ResearchBriefID string
	GuidanceKey     string
	GuidanceValue   string
}

func CreateAgentGuidance(
	ctx context.Context,
	dbtx db.DBTX,
	data CreateAgentGuidanceData,
) (AgentGuidance, error) {
	if err := validate.Struct(data); err != nil {
		return AgentGuidance{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.NewInsertAgentGuidanceParams(
		data.ResearchBriefID,
		data.GuidanceKey,
		data.GuidanceValue,
	)
	row, err := db.New().InsertAgentGuidance(ctx, dbtx, params)
	if err != nil {
		return AgentGuidance{}, err
	}

	result, err := rowToAgentGuidance(row)
	if err != nil {
		return AgentGuidance{}, err
	}
	return result, nil
}

type UpdateAgentGuidanceData struct {
	ID              uuid.UUID
	ResearchBriefID string
	GuidanceKey     string
	GuidanceValue   string
}

func UpdateAgentGuidance(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateAgentGuidanceData,
) (AgentGuidance, error) {
	if err := validate.Struct(data); err != nil {
		return AgentGuidance{}, errors.Join(ErrDomainValidation, err)
	}

	currentRow, err := db.New().QueryAgentGuidanceByID(ctx, dbtx, data.ID.String())
	if err != nil {
		return AgentGuidance{}, err
	}

	params := db.NewUpdateAgentGuidanceParams(
		data.ID.String(),
		func() string {
			if true {
				return data.ResearchBriefID
			}
			return currentRow.ResearchBriefID
		}(),
		func() string {
			if true {
				return data.GuidanceKey
			}
			return currentRow.GuidanceKey
		}(),
		func() string {
			if true {
				return data.GuidanceValue
			}
			return currentRow.GuidanceValue
		}(),
	)

	row, err := db.New().UpdateAgentGuidance(ctx, dbtx, params)
	if err != nil {
		return AgentGuidance{}, err
	}

	result, err := rowToAgentGuidance(row)
	if err != nil {
		return AgentGuidance{}, err
	}
	return result, nil
}

func DestroyAgentGuidance(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.New().DeleteAgentGuidance(ctx, dbtx, id.String())
}

func AllAgentGuidances(
	ctx context.Context,
	dbtx db.DBTX,
) ([]AgentGuidance, error) {
	rows, err := db.New().QueryAllAgentGuidances(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	agentguidances := make([]AgentGuidance, len(rows))
	for i, row := range rows {
		result, err := rowToAgentGuidance(row)
		if err != nil {
			return nil, err
		}
		agentguidances[i] = result
	}

	return agentguidances, nil
}

type PaginatedAgentGuidances struct {
	AgentGuidances []AgentGuidance
	TotalCount     int64
	Page           int64
	PageSize       int64
	TotalPages     int64
}

func PaginateAgentGuidances(
	ctx context.Context,
	dbtx db.DBTX,
	page int64,
	pageSize int64,
) (PaginatedAgentGuidances, error) {
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

	totalCount, err := db.New().CountAgentGuidances(ctx, dbtx)
	if err != nil {
		return PaginatedAgentGuidances{}, err
	}

	rows, err := db.New().QueryPaginatedAgentGuidances(
		ctx,
		dbtx,
		db.NewQueryPaginatedAgentGuidancesParams(pageSize, offset),
	)
	if err != nil {
		return PaginatedAgentGuidances{}, err
	}

	agentguidances := make([]AgentGuidance, len(rows))
	for i, row := range rows {
		result, err := rowToAgentGuidance(row)
		if err != nil {
			return PaginatedAgentGuidances{}, err
		}
		agentguidances[i] = result
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedAgentGuidances{
		AgentGuidances: agentguidances,
		TotalCount:     totalCount,
		Page:           page,
		PageSize:       pageSize,
		TotalPages:     totalPages,
	}, nil
}

func rowToAgentGuidance(row db.Agentguidance) (AgentGuidance, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return AgentGuidance{}, err
	}

	return AgentGuidance{
		ID:              id,
		ResearchBriefID: row.ResearchBriefID,
		GuidanceKey:     row.GuidanceKey,
		GuidanceValue:   row.GuidanceValue,
	}, nil
}
