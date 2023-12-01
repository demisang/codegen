package db

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const __camelCase__TableName = "__snake_case__"

type __PascalCase__Repository struct {
	conn *pgxpool.Pool
	log  *logrus.Entry
}

func New__PascalCase__Repository(db *Connection) *__PascalCase__Repository {
	return &__PascalCase__Repository{
		conn: db.pool,
		log:  db.log.WithField("module", "__camelCase__Repository"),
	}
}

func (r *__PascalCase__Repository) baseQuery() sq.SelectBuilder {
	return sq.
		Select(
			"id",
			"name",
			// @TODO database fields here
		).
		From(__camelCase__TableName).
		PlaceholderFormat(sq.Dollar)
}

func (r *__PascalCase__Repository) baseQueryScan(row pgx.Row, model *models.__PascalCase__) error {
	return row.Scan(
		&model.ID,
		&model.Name,
		// @TODO database fields here
	)
}

func (r *__PascalCase__Repository) GetByID(ctx context.Context, id models.__PascalCase__ID) (m models.__PascalCase__, err error) {
	query := r.baseQuery().
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return m, fmt.Errorf("sql builder query row: %w", err)
	}

	row := r.conn.QueryRow(ctx, sql, args...)
	err = r.baseQueryScan(row, &m)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return m, models.Err__PascalCase__NotFound
	case err != nil:
		return m, fmt.Errorf("query row scan: %w", err)
	}

	return m, nil
}

func (r *__PascalCase__Repository) GetAll(ctx context.Context, filter models.__PascalCase__Filter) (*[]models.__PascalCase__, error) {
	query := r.baseQuery().
		Limit(filter.Limit).
		Offset(filter.Offset)

	// Filtering statements
	if filter.IDs != nil {
		query = query.Where(sq.Eq{"id": *filter.IDs})
	}

	// @TODO database fields here

	if filter.SearchString != nil {
		query = query.Where(sq.Like{"name": *filter.SearchString + "%"})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql builder query all: %w", err)
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query all: %w", err)
	}

	// Models list
	results := make([]models.__PascalCase__, 0)

	for rows.Next() {
		var model models.__PascalCase__

		err = r.baseQueryScan(rows, &model)
		if err != nil {
			return nil, fmt.Errorf("query all scan: %w", err)
		}

		results = append(results, model)
	}

	return &results, nil
}

func (r *__PascalCase__Repository) Create(ctx context.Context, dto models.__PascalCase__CreateDto) (models.__PascalCase__ID, error) {
	query := sq.
		Insert(__camelCase__TableName).
		Columns(
			"name",
			// @TODO database fields here
		).
		Values(
			dto.Name,
			// @TODO database fields here
		).
		Suffix(`RETURNING "id"`)

	var id models.__PascalCase__ID

	sql, args, err := query.ToSql()
	if err != nil {
		return id, fmt.Errorf("sql builder insert: %w", err)
	}

	err = r.conn.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("insert: %w", err)
	}

	return id, nil
}

func (r *__PascalCase__Repository) Update(ctx context.Context, id models.__PascalCase__ID, dto models.__PascalCase__UpdateDto) error {
	query := sq.
		Update(__camelCase__TableName).
		SetMap(map[string]interface{}{
			"name": dto.Name,
			// @TODO database fields here
		}).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("sql builder update: %w", err)
	}

	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (r *__PascalCase__Repository) Delete(ctx context.Context, id models.__PascalCase__ID) error {
	query := sq.
		Delete(__camelCase__TableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("sql builder delete: %w", err)
	}

	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
