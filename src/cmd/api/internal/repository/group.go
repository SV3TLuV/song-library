package repository

import (
	"context"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"song-library-api/src/cmd/api/internal/converter"
	"song-library-api/src/cmd/api/internal/model"
)

var _ GroupRepository = (*groupRepository)(nil)

type groupRepository struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	trManager *manager.Manager
}

func NewGroupRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trManager *manager.Manager) *groupRepository {
	return &groupRepository{
		pool:      pool,
		getter:    getter,
		trManager: trManager,
	}
}

func (repo *groupRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error) {
	query := goqu.Dialect("postgres").
		From("group").
		Where(goqu.Ex{"id": id})

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var group model.Group
	tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
	if err = pgxscan.Get(ctx, tr, &group, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(model.ErrNotFound, "group not found")
		}
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return &group, nil
}

func (repo *groupRepository) GetByName(ctx context.Context, name string) (*model.Group, error) {
	query := goqu.Dialect("postgres").
		From("group").
		Where(goqu.Ex{"name": name})

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var group model.Group
	tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
	if err = pgxscan.Get(ctx, tr, &group, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(model.ErrNotFound, "group not found")
		}
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return &group, nil
}

func (repo *groupRepository) Create(ctx context.Context, entity model.Group) (*model.Group, error) {
	var group model.Group
	err := repo.trManager.Do(ctx, func(ctx context.Context) error {
		tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
		record := converter.ToRecordFromGroup(entity)

		query := goqu.Dialect("postgres").
			Insert("group").
			Rows(record).
			Returning("group.*")

		sql, args, err := query.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		if err = pgxscan.Get(ctx, tr, &group, sql, args...); err != nil {
			return errors.Wrap(err, "failed to execute query")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return &group, nil
}

func (repo *groupRepository) Update(ctx context.Context, entity model.Group) (*model.Group, error) {
	var group model.Group
	err := repo.trManager.Do(ctx, func(ctx context.Context) error {
		tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
		record := converter.ToRecordFromGroup(entity)
		record["updated_at"] = goqu.L("CURRENT_TIMESTAMP")

		query := goqu.Dialect("postgres").
			Update("group").
			Set(record).
			Where(goqu.Ex{"id": entity.ID}).
			Returning("group.*")

		sql, args, err := query.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		if err = pgxscan.Get(ctx, tr, &group, sql, args...); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrap(model.ErrNotFound, "group not found")
			}
			return errors.Wrap(err, "failed to execute query")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return &group, nil
}
