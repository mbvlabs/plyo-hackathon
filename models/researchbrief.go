package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/mbvlabs/plyo-hackathon/models/internal/db"
)

type ResearchBrief struct {
	ID                   uuid.UUID
	IdentificationStatus string
	CompanyName          string
	OfficialDomain       string
	Headquarters         string
	Industry             string
	CompanyType          string
	Status               string
	GeographicScope      string
	ResearchDepth        string
	ConfidenceScore      float64
	LastUpdated          time.Time
}

func FindResearchBrief(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (ResearchBrief, error) {
	row, err := db.New().QueryResearchBriefByID(ctx, dbtx, id.String())
	if err != nil {
		return ResearchBrief{}, err
	}

	result, err := rowToResearchBrief(row)
	if err != nil {
		return ResearchBrief{}, err
	}
	return result, nil
}

type CreateResearchBriefData struct {
	IdentificationStatus string
	CompanyName          string
	OfficialDomain       string
	Headquarters         string
	Industry             string
	CompanyType          string
	Status               string
	GeographicScope      string
	ResearchDepth        string
	ConfidenceScore      float64
	LastUpdated          time.Time
}

func CreateResearchBrief(
	ctx context.Context,
	dbtx db.DBTX,
	data CreateResearchBriefData,
) (ResearchBrief, error) {
	if err := validate.Struct(data); err != nil {
		return ResearchBrief{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.NewInsertResearchBriefParams(
		data.IdentificationStatus,
		data.CompanyName,
		data.OfficialDomain,
		data.Headquarters,
		data.Industry,
		data.CompanyType,
		data.Status,
		data.GeographicScope,
		data.ResearchDepth,
		data.ConfidenceScore,
		data.LastUpdated,
	)
	row, err := db.New().InsertResearchBrief(ctx, dbtx, params)
	if err != nil {
		return ResearchBrief{}, err
	}

	result, err := rowToResearchBrief(row)
	if err != nil {
		return ResearchBrief{}, err
	}
	return result, nil
}

type UpdateResearchBriefData struct {
	ID                   uuid.UUID
	IdentificationStatus string
	CompanyName          string
	OfficialDomain       string
	Headquarters         string
	Industry             string
	CompanyType          string
	Status               string
	GeographicScope      string
	ResearchDepth        string
	ConfidenceScore      float64
	LastUpdated          time.Time
}

func UpdateResearchBrief(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateResearchBriefData,
) (ResearchBrief, error) {
	if err := validate.Struct(data); err != nil {
		return ResearchBrief{}, errors.Join(ErrDomainValidation, err)
	}

	currentRow, err := db.New().QueryResearchBriefByID(ctx, dbtx, data.ID.String())
	if err != nil {
		return ResearchBrief{}, err
	}

	params := db.NewUpdateResearchBriefParams(
		data.ID.String(),
		func() string {
			if true {
				return data.IdentificationStatus
			}
			return currentRow.IdentificationStatus
		}(),
		func() string {
			if true {
				return data.CompanyName
			}
			return currentRow.CompanyName
		}(),
		func() string {
			if true {
				return data.OfficialDomain
			}
			return currentRow.OfficialDomain
		}(),
		func() string {
			if true {
				return data.Headquarters
			}
			return currentRow.Headquarters
		}(),
		func() string {
			if true {
				return data.Industry
			}
			return currentRow.Industry
		}(),
		func() string {
			if true {
				return data.CompanyType
			}
			return currentRow.CompanyType
		}(),
		func() string {
			if true {
				return data.Status
			}
			return currentRow.Status
		}(),
		func() string {
			if true {
				return data.GeographicScope
			}
			return currentRow.GeographicScope
		}(),
		func() string {
			if true {
				return data.ResearchDepth
			}
			return currentRow.ResearchDepth
		}(),
		func() float64 {
			if true {
				return data.ConfidenceScore
			}
			return currentRow.ConfidenceScore
		}(),
		func() time.Time {
			if true {
				return data.LastUpdated
			}
			return currentRow.LastUpdated
		}(),
	)

	row, err := db.New().UpdateResearchBrief(ctx, dbtx, params)
	if err != nil {
		return ResearchBrief{}, err
	}

	result, err := rowToResearchBrief(row)
	if err != nil {
		return ResearchBrief{}, err
	}
	return result, nil
}

func DestroyResearchBrief(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.New().DeleteResearchBrief(ctx, dbtx, id.String())
}

func AllResearchBriefs(
	ctx context.Context,
	dbtx db.DBTX,
) ([]ResearchBrief, error) {
	rows, err := db.New().QueryAllResearchBriefs(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	researchbriefs := make([]ResearchBrief, len(rows))
	for i, row := range rows {
		result, err := rowToResearchBrief(row)
		if err != nil {
			return nil, err
		}
		researchbriefs[i] = result
	}

	return researchbriefs, nil
}

type PaginatedResearchBriefs struct {
	ResearchBriefs []ResearchBrief
	TotalCount     int64
	Page           int64
	PageSize       int64
	TotalPages     int64
}

func PaginateResearchBriefs(
	ctx context.Context,
	dbtx db.DBTX,
	page int64,
	pageSize int64,
) (PaginatedResearchBriefs, error) {
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

	totalCount, err := db.New().CountResearchBriefs(ctx, dbtx)
	if err != nil {
		return PaginatedResearchBriefs{}, err
	}

	rows, err := db.New().QueryPaginatedResearchBriefs(
		ctx,
		dbtx,
		db.NewQueryPaginatedResearchBriefsParams(pageSize, offset),
	)
	if err != nil {
		return PaginatedResearchBriefs{}, err
	}

	researchbriefs := make([]ResearchBrief, len(rows))
	for i, row := range rows {
		result, err := rowToResearchBrief(row)
		if err != nil {
			return PaginatedResearchBriefs{}, err
		}
		researchbriefs[i] = result
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedResearchBriefs{
		ResearchBriefs: researchbriefs,
		TotalCount:     totalCount,
		Page:           page,
		PageSize:       pageSize,
		TotalPages:     totalPages,
	}, nil
}

func rowToResearchBrief(row db.Researchbrief) (ResearchBrief, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return ResearchBrief{}, err
	}

	return ResearchBrief{
		ID:                   id,
		IdentificationStatus: row.IdentificationStatus,
		CompanyName:          row.CompanyName,
		OfficialDomain:       row.OfficialDomain,
		Headquarters:         row.Headquarters,
		Industry:             row.Industry,
		CompanyType:          row.CompanyType,
		Status:               row.Status,
		GeographicScope:      row.GeographicScope,
		ResearchDepth:        row.ResearchDepth,
		ConfidenceScore:      row.ConfidenceScore,
		LastUpdated:          row.LastUpdated,
	}, nil
}
