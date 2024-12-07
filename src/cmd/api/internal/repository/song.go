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

var _ SongRepository = (*songRepository)(nil)

type songRepository struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	trManager *manager.Manager
}

func NewSongRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trManager *manager.Manager) *songRepository {
	return &songRepository{
		pool:      pool,
		getter:    getter,
		trManager: trManager,
	}
}

func (repo *songRepository) GetSongs(ctx context.Context,
	filters *model.SongFilter,
	limit, offset uint) ([]model.Song, error) {
	query := goqu.Dialect("postgres").
		From("song").
		Join(
			goqu.T("group"),
			goqu.On(goqu.I("song.group_id").Eq(goqu.I("group.id"))),
		).
		Select(
			goqu.I("song.*"),
			goqu.I("group.name").As("group"),
		).
		Order(goqu.I("id").Asc()).
		Limit(limit).
		Offset(offset)

	if filters != nil {
		if filters.GroupID != uuid.Nil {
			query = query.Where(goqu.Ex{"group_id": filters.GroupID})
		}
		if filters.Song != nil {
			query = query.Where(goqu.Ex{"song": goqu.Op{"ilike": "%" + *filters.Song + "%"}})
		}
		if filters.Text != nil {
			query = query.Where(goqu.Ex{"text": goqu.Op{"ilike": "%" + *filters.Text + "%"}})
		}
		if filters.Link != nil {
			query = query.Where(goqu.Ex{"link": goqu.Op{"ilike": "%" + *filters.Link + "%"}})
		}
		if filters.ReleaseDate != nil {
			query = query.Where(goqu.I("release_date").Eq(*filters.ReleaseDate))
		}
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	songs := make([]model.Song, 0)
	tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
	err = pgxscan.Select(ctx, tr, &songs, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return songs, nil
}

func (repo *songRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Song, error) {
	query := goqu.Dialect("postgres").
		From("song").
		Join(
			goqu.T("group"),
			goqu.On(goqu.I("song.group_id").Eq(goqu.I("group.id"))),
		).
		Select(
			goqu.I("song.*"),
			goqu.I("group.name").As("group"),
		).
		Where(goqu.Ex{"song.id": id})

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var song model.Song
	tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
	if err = pgxscan.Get(ctx, tr, &song, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(model.ErrNotFound, "song not found")
		}
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return &song, nil
}

func (repo *songRepository) GetByNameAndGroup(ctx context.Context, group, name string) (*model.Song, error) {
	query := goqu.Dialect("postgres").
		From("song").
		Join(
			goqu.T("group"),
			goqu.On(goqu.I("song.group_id").Eq(goqu.I("group.id"))),
		).
		Select(
			goqu.I("song.*"),
			goqu.I("group.name").As("group"),
		).
		Where(goqu.And(goqu.Ex{"name": name}, goqu.Ex{"group": group}))

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var song model.Song
	tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
	if err = pgxscan.Get(ctx, tr, &song, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(model.ErrNotFound, "song not found")
		}
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return &song, nil
}

func (repo *songRepository) Count(ctx context.Context, filters *model.SongFilter) (uint, error) {
	query := goqu.Dialect("postgres").
		From("song").
		Select(goqu.COUNT("*").As("total"))

	if filters != nil {
		if filters.GroupID != uuid.Nil {
			query = query.Where(goqu.Ex{"group_id": filters.GroupID})
		}
		if filters.Song != nil {
			query = query.Where(goqu.Ex{"song": goqu.Op{"ilike": "%" + *filters.Song + "%"}})
		}
		if filters.Text != nil {
			query = query.Where(goqu.Ex{"text": goqu.Op{"ilike": "%" + *filters.Text + "%"}})
		}
		if filters.Link != nil {
			query = query.Where(goqu.Ex{"link": goqu.Op{"ilike": "%" + *filters.Link + "%"}})
		}
		if filters.ReleaseDate != nil {
			query = query.Where(goqu.I("release_date").Eq(*filters.ReleaseDate))
		}
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build query")
	}

	var total uint
	tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
	err = tr.QueryRow(ctx, sql, args...).Scan(&total)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute query")
	}

	return total, nil
}

func (repo *songRepository) Create(ctx context.Context, entity model.Song) (*model.Song, error) {
	var song model.Song
	err := repo.trManager.Do(ctx, func(ctx context.Context) error {
		tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
		record := converter.ToRecordFromSong(entity)

		query := goqu.Dialect("postgres").
			Insert("song").
			Rows(record).
			Returning("song.*")

		sql, args, err := query.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		if err = pgxscan.Get(ctx, tr, &song, sql, args...); err != nil {
			return errors.Wrap(err, "failed to execute query")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return &song, nil
}

func (repo *songRepository) Update(ctx context.Context, entity model.Song) (*model.Song, error) {
	var song model.Song
	err := repo.trManager.Do(ctx, func(ctx context.Context) error {
		tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
		record := converter.ToRecordFromSong(entity)
		record["updated_at"] = goqu.L("CURRENT_TIMESTAMP")

		query := goqu.Dialect("postgres").
			Update("song").
			Set(record).
			Where(goqu.Ex{"id": entity.ID}).
			Returning("song.*")

		sql, args, err := query.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		if err = pgxscan.Get(ctx, tr, &song, sql, args...); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrap(model.ErrNotFound, "song not found")
			}
			return errors.Wrap(err, "failed to execute query")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return &song, nil
}

func (repo *songRepository) Delete(ctx context.Context, song model.Song) error {
	query := goqu.Dialect("postgres").
		Delete("song").
		Where(goqu.Ex{
			"id": song.ID,
		})

	sql, args, err := query.ToSQL()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	tr := repo.getter.DefaultTrOrDB(ctx, repo.pool)
	if _, err = tr.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to execute query")
	}

	return nil
}
