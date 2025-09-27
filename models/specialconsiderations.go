package models

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/mbvlabs/plyo-hackathon/models/internal/db"
)

type SpecialConsiderations struct {
	ID              uuid.UUID
	ResearchBriefID string
	Consideration   string
}

func FindSpecialConsiderations(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (SpecialConsiderations, error) {
	row, err := db.New().QuerySpecialConsiderationsByID(ctx, dbtx, id.String())
	if err != nil {
		return SpecialConsiderations{}, err
	}

	result, err := rowToSpecialConsiderations(row)
	if err != nil {
		return SpecialConsiderations{}, err
	}
	return result, nil
}

func FindSpecialConsiderationsByResearchBriefID(
	ctx context.Context,
	dbtx db.DBTX,
	researchBriefID string,
) ([]SpecialConsiderations, error) {
	rows, err := db.New().QuerySpecialConsiderationsByResearchBriefID(ctx, dbtx, researchBriefID)
	if err != nil {
		return nil, err
	}

	considerations := make([]SpecialConsiderations, len(rows))
	for i, row := range rows {
		result, err := rowToSpecialConsiderations(row)
		if err != nil {
			return nil, err
		}
		considerations[i] = result
	}

	return considerations, nil
}

type CreateSpecialConsiderationsData struct {
	ResearchBriefID string
	Consideration   string
}

func CreateSpecialConsiderations(
	ctx context.Context,
	dbtx db.DBTX,
	data CreateSpecialConsiderationsData,
) (SpecialConsiderations, error) {
	if err := validate.Struct(data); err != nil {
		return SpecialConsiderations{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.NewInsertSpecialConsiderationsParams(
		data.ResearchBriefID,
		data.Consideration,
	)
	row, err := db.New().InsertSpecialConsiderations(ctx, dbtx, params)
	if err != nil {
		return SpecialConsiderations{}, err
	}

	result, err := rowToSpecialConsiderations(row)
	if err != nil {
		return SpecialConsiderations{}, err
	}
	return result, nil
}

type UpdateSpecialConsiderationsData struct {
	ID              uuid.UUID
	ResearchBriefID string
	Consideration   string
}

func UpdateSpecialConsiderations(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateSpecialConsiderationsData,
) (SpecialConsiderations, error) {
	if err := validate.Struct(data); err != nil {
		return SpecialConsiderations{}, errors.Join(ErrDomainValidation, err)
	}

	currentRow, err := db.New().QuerySpecialConsiderationsByID(ctx, dbtx, data.ID.String())
	if err != nil {
		return SpecialConsiderations{}, err
	}

	params := db.NewUpdateSpecialConsiderationsParams(
		data.ID.String(),
		func() string {
			if true {
				return data.ResearchBriefID
			}
			return currentRow.ResearchBriefID
		}(),
		func() string {
			if true {
				return data.Consideration
			}
			return currentRow.Consideration
		}(),
	)

	row, err := db.New().UpdateSpecialConsiderations(ctx, dbtx, params)
	if err != nil {
		return SpecialConsiderations{}, err
	}

	result, err := rowToSpecialConsiderations(row)
	if err != nil {
		return SpecialConsiderations{}, err
	}
	return result, nil
}

func DestroySpecialConsiderations(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.New().DeleteSpecialConsiderations(ctx, dbtx, id.String())
}

func AllSpecialConsiderationss(
	ctx context.Context,
	dbtx db.DBTX,
) ([]SpecialConsiderations, error) {
	rows, err := db.New().QueryAllSpecialConsiderationss(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	specialconsiderationss := make([]SpecialConsiderations, len(rows))
	for i, row := range rows {
		result, err := rowToSpecialConsiderations(row)
		if err != nil {
			return nil, err
		}
		specialconsiderationss[i] = result
	}

	return specialconsiderationss, nil
}

type PaginatedSpecialConsiderationss struct {
	SpecialConsiderationss []SpecialConsiderations
	TotalCount             int64
	Page                   int64
	PageSize               int64
	TotalPages             int64
}

func PaginateSpecialConsiderationss(
	ctx context.Context,
	dbtx db.DBTX,
	page int64,
	pageSize int64,
) (PaginatedSpecialConsiderationss, error) {
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

	totalCount, err := db.New().CountSpecialConsiderationss(ctx, dbtx)
	if err != nil {
		return PaginatedSpecialConsiderationss{}, err
	}

	rows, err := db.New().QueryPaginatedSpecialConsiderationss(
		ctx,
		dbtx,
		db.NewQueryPaginatedSpecialConsiderationssParams(pageSize, offset),
	)
	if err != nil {
		return PaginatedSpecialConsiderationss{}, err
	}

	specialconsiderationss := make([]SpecialConsiderations, len(rows))
	for i, row := range rows {
		result, err := rowToSpecialConsiderations(row)
		if err != nil {
			return PaginatedSpecialConsiderationss{}, err
		}
		specialconsiderationss[i] = result
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedSpecialConsiderationss{
		SpecialConsiderationss: specialconsiderationss,
		TotalCount:             totalCount,
		Page:                   page,
		PageSize:               pageSize,
		TotalPages:             totalPages,
	}, nil
}

func rowToSpecialConsiderations(row db.Specialconsideration) (SpecialConsiderations, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return SpecialConsiderations{}, err
	}

	return SpecialConsiderations{
		ID:              id,
		ResearchBriefID: row.ResearchBriefID,
		Consideration:   row.Consideration,
	}, nil
}
