package models

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/mbvlabs/plyo-hackathon/models/internal/db"
)

type CompanyCandidates struct {
	ID              uuid.UUID
	ResearchBriefID string
	Name            string
	Domain          string
	Description     string
	Industry        string
	Location        string
}

func FindCompanyCandidates(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (CompanyCandidates, error) {
	row, err := db.New().QueryCompanyCandidatesByID(ctx, dbtx, id.String())
	if err != nil {
		return CompanyCandidates{}, err
	}

	result, err := rowToCompanyCandidates(row)
	if err != nil {
		return CompanyCandidates{}, err
	}
	return result, nil
}

func FindCompanyCandidatesByResearchBriefID(
	ctx context.Context,
	dbtx db.DBTX,
	researchBriefID string,
) ([]CompanyCandidates, error) {
	rows, err := db.New().QueryCompanyCandidatesByResearchBriefID(ctx, dbtx, researchBriefID)
	if err != nil {
		return nil, err
	}

	candidates := make([]CompanyCandidates, len(rows))
	for i, row := range rows {
		result, err := rowToCompanyCandidates(row)
		if err != nil {
			return nil, err
		}
		candidates[i] = result
	}

	return candidates, nil
}

type CreateCompanyCandidatesData struct {
	ResearchBriefID string
	Name            string
	Domain          string
	Description     string
	Industry        string
	Location        string
}

func CreateCompanyCandidates(
	ctx context.Context,
	dbtx db.DBTX,
	data CreateCompanyCandidatesData,
) (CompanyCandidates, error) {
	if err := validate.Struct(data); err != nil {
		return CompanyCandidates{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.NewInsertCompanyCandidatesParams(
		data.ResearchBriefID,
		data.Name,
		data.Domain,
		data.Description,
		data.Industry,
		data.Location,
	)
	row, err := db.New().InsertCompanyCandidates(ctx, dbtx, params)
	if err != nil {
		return CompanyCandidates{}, err
	}

	result, err := rowToCompanyCandidates(row)
	if err != nil {
		return CompanyCandidates{}, err
	}
	return result, nil
}

type UpdateCompanyCandidatesData struct {
	ID              uuid.UUID
	ResearchBriefID string
	Name            string
	Domain          string
	Description     string
	Industry        string
	Location        string
}

func UpdateCompanyCandidates(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateCompanyCandidatesData,
) (CompanyCandidates, error) {
	if err := validate.Struct(data); err != nil {
		return CompanyCandidates{}, errors.Join(ErrDomainValidation, err)
	}

	currentRow, err := db.New().QueryCompanyCandidatesByID(ctx, dbtx, data.ID.String())
	if err != nil {
		return CompanyCandidates{}, err
	}

	params := db.NewUpdateCompanyCandidatesParams(
		data.ID.String(),
		func() string {
			if true {
				return data.ResearchBriefID
			}
			return currentRow.ResearchBriefID
		}(),
		func() string {
			if true {
				return data.Name
			}
			return currentRow.Name
		}(),
		func() string {
			if true {
				return data.Domain
			}
			return currentRow.Domain
		}(),
		func() string {
			if true {
				return data.Description
			}
			return currentRow.Description
		}(),
		func() string {
			if true {
				return data.Industry
			}
			return currentRow.Industry
		}(),
		func() string {
			if true {
				return data.Location
			}
			return currentRow.Location
		}(),
	)

	row, err := db.New().UpdateCompanyCandidates(ctx, dbtx, params)
	if err != nil {
		return CompanyCandidates{}, err
	}

	result, err := rowToCompanyCandidates(row)
	if err != nil {
		return CompanyCandidates{}, err
	}
	return result, nil
}

func DestroyCompanyCandidates(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.New().DeleteCompanyCandidates(ctx, dbtx, id.String())
}

func AllCompanyCandidatess(
	ctx context.Context,
	dbtx db.DBTX,
) ([]CompanyCandidates, error) {
	rows, err := db.New().QueryAllCompanyCandidatess(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	companycandidatess := make([]CompanyCandidates, len(rows))
	for i, row := range rows {
		result, err := rowToCompanyCandidates(row)
		if err != nil {
			return nil, err
		}
		companycandidatess[i] = result
	}

	return companycandidatess, nil
}

type PaginatedCompanyCandidatess struct {
	CompanyCandidatess []CompanyCandidates
	TotalCount         int64
	Page               int64
	PageSize           int64
	TotalPages         int64
}

func PaginateCompanyCandidatess(
	ctx context.Context,
	dbtx db.DBTX,
	page int64,
	pageSize int64,
) (PaginatedCompanyCandidatess, error) {
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

	totalCount, err := db.New().CountCompanyCandidatess(ctx, dbtx)
	if err != nil {
		return PaginatedCompanyCandidatess{}, err
	}

	rows, err := db.New().QueryPaginatedCompanyCandidatess(
		ctx,
		dbtx,
		db.NewQueryPaginatedCompanyCandidatessParams(pageSize, offset),
	)
	if err != nil {
		return PaginatedCompanyCandidatess{}, err
	}

	companycandidatess := make([]CompanyCandidates, len(rows))
	for i, row := range rows {
		result, err := rowToCompanyCandidates(row)
		if err != nil {
			return PaginatedCompanyCandidatess{}, err
		}
		companycandidatess[i] = result
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedCompanyCandidatess{
		CompanyCandidatess: companycandidatess,
		TotalCount:         totalCount,
		Page:               page,
		PageSize:           pageSize,
		TotalPages:         totalPages,
	}, nil
}

func rowToCompanyCandidates(row db.Companycandidate) (CompanyCandidates, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return CompanyCandidates{}, err
	}

	return CompanyCandidates{
		ID:              id,
		ResearchBriefID: row.ResearchBriefID,
		Name:            row.Name,
		Domain:          row.Domain,
		Description:     row.Description,
		Industry:        row.Industry,
		Location:        row.Location,
	}, nil
}
