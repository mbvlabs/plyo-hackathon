package models

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/mbvlabs/plyo-hackathon/models/internal/db"
)

type Sources struct {
	ID              uuid.UUID
	ResearchBriefID string
	SourceUrl       string
}

func FindSources(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Sources, error) {
	row, err := db.New().QuerySourcesByID(ctx, dbtx, id.String())
	if err != nil {
		return Sources{}, err
	}

	result, err := rowToSources(row)
	if err != nil {
		return Sources{}, err
	}
	return result, nil
}

func FindSourcesByResearchBriefID(
	ctx context.Context,
	dbtx db.DBTX,
	researchBriefID string,
) ([]Sources, error) {
	rows, err := db.New().QuerySourcesByResearchBriefID(ctx, dbtx, researchBriefID)
	if err != nil {
		return nil, err
	}

	sources := make([]Sources, len(rows))
	for i, row := range rows {
		result, err := rowToSources(row)
		if err != nil {
			return nil, err
		}
		sources[i] = result
	}

	return sources, nil
}

type CreateSourcesData struct {
	ResearchBriefID string
	SourceUrl       string
}

func CreateSources(
	ctx context.Context,
	dbtx db.DBTX,
	data CreateSourcesData,
) (Sources, error) {
	if err := validate.Struct(data); err != nil {
		return Sources{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.NewInsertSourcesParams(
		data.ResearchBriefID,
		data.SourceUrl,
	)
	row, err := db.New().InsertSources(ctx, dbtx, params)
	if err != nil {
		return Sources{}, err
	}

	result, err := rowToSources(row)
	if err != nil {
		return Sources{}, err
	}
	return result, nil
}

type UpdateSourcesData struct {
	ID              uuid.UUID
	ResearchBriefID string
	SourceUrl       string
}

func UpdateSources(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateSourcesData,
) (Sources, error) {
	if err := validate.Struct(data); err != nil {
		return Sources{}, errors.Join(ErrDomainValidation, err)
	}

	currentRow, err := db.New().QuerySourcesByID(ctx, dbtx, data.ID.String())
	if err != nil {
		return Sources{}, err
	}

	params := db.NewUpdateSourcesParams(
		data.ID.String(),
		func() string {
			if true {
				return data.ResearchBriefID
			}
			return currentRow.ResearchBriefID
		}(),
		func() string {
			if true {
				return data.SourceUrl
			}
			return currentRow.SourceUrl
		}(),
	)

	row, err := db.New().UpdateSources(ctx, dbtx, params)
	if err != nil {
		return Sources{}, err
	}

	result, err := rowToSources(row)
	if err != nil {
		return Sources{}, err
	}
	return result, nil
}

func DestroySources(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.New().DeleteSources(ctx, dbtx, id.String())
}

func AllSourcess(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Sources, error) {
	rows, err := db.New().QueryAllSourcess(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	sourcess := make([]Sources, len(rows))
	for i, row := range rows {
		result, err := rowToSources(row)
		if err != nil {
			return nil, err
		}
		sourcess[i] = result
	}

	return sourcess, nil
}

type PaginatedSourcess struct {
	Sourcess   []Sources
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func PaginateSourcess(
	ctx context.Context,
	dbtx db.DBTX,
	page int64,
	pageSize int64,
) (PaginatedSourcess, error) {
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

	totalCount, err := db.New().CountSourcess(ctx, dbtx)
	if err != nil {
		return PaginatedSourcess{}, err
	}

	rows, err := db.New().QueryPaginatedSourcess(
		ctx,
		dbtx,
		db.NewQueryPaginatedSourcessParams(pageSize, offset),
	)
	if err != nil {
		return PaginatedSourcess{}, err
	}

	sourcess := make([]Sources, len(rows))
	for i, row := range rows {
		result, err := rowToSources(row)
		if err != nil {
			return PaginatedSourcess{}, err
		}
		sourcess[i] = result
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedSourcess{
		Sourcess:   sourcess,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func rowToSources(row db.Source) (Sources, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return Sources{}, err
	}

	return Sources{
		ID:              id,
		ResearchBriefID: row.ResearchBriefID,
		SourceUrl:       row.SourceUrl,
	}, nil
}
