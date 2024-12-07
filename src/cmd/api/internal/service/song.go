package service

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log/slog"
	"math"
	"song-library-api/src/cmd/api/internal/model"
	"song-library-api/src/cmd/api/internal/repository"
	"song-library-api/src/pkg/music_info_client"
	"strings"
	"time"
)

var _ SongService = (*songService)(nil)

type songService struct {
	songRepo        repository.SongRepository
	groupRepo       repository.GroupRepository
	musicInfoClient *music_info_client.MusicInfoClient
	trManager       *manager.Manager
	logger          *slog.Logger
}

func NewSongService(
	songRepo repository.SongRepository,
	groupRepo repository.GroupRepository,
	musicInfoClient *music_info_client.MusicInfoClient,
	trManager *manager.Manager,
	logger *slog.Logger) *songService {
	return &songService{
		songRepo:        songRepo,
		groupRepo:       groupRepo,
		musicInfoClient: musicInfoClient,
		trManager:       trManager,
		logger:          logger,
	}
}

func (s *songService) GetSongs(ctx context.Context, filters *model.SongFilter, page, pageSize uint) (*model.PaginatedList[model.Song], error) {
	offset := (page - 1) * pageSize
	songs, err := s.songRepo.GetSongs(ctx, filters, pageSize, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.songRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	totalPages := uint(math.Ceil(float64(total) / float64(pageSize)))

	s.logger.Info("get songs", "filters", filters)

	return &model.PaginatedList[model.Song]{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Items:      songs,
	}, nil
}

func (s *songService) GetSongText(ctx context.Context, id uuid.UUID, page, pageSize uint) (*model.PaginatedList[string], error) {
	song, err := s.songRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get song failed")
	}

	verses := strings.Split(song.Text, "\n\n")
	total := uint(len(verses))
	totalPages := uint(math.Ceil(float64(total) / float64(pageSize)))

	if page > totalPages {
		return nil, errors.Wrap(model.ErrBadRequest, "page out of range")
	}

	start := (page - 1) * pageSize
	if start > total {
		return nil, errors.Wrap(model.ErrBadRequest, "page out of range")
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	s.logger.Info("get song text", "id", id, "song", song, "group", song.Group)

	return &model.PaginatedList[string]{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Items:      verses[start:end],
	}, nil
}

func (s *songService) GetByID(ctx context.Context, id uuid.UUID) (*model.Song, error) {
	song, err := s.songRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.logger.Info("get song", "id", id, "song", song, "group", song.Group)

	return song, nil
}

func (s *songService) Add(ctx context.Context, song, group string) (*model.Song, error) {
	songDetail, err := s.musicInfoClient.GetSongInfo(ctx, group, song)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get song detail")
	}

	var created *model.Song
	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		groupDB, err := s.groupRepo.GetByName(ctx, group)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return errors.Wrap(err, "failed to get group by name")
		}

		if groupDB == nil {
			groupDB, err = s.groupRepo.Create(ctx, model.Group{Name: group})
			if err != nil {
				return err
			}

			s.logger.Info("group created", "id", groupDB.ID, "group", groupDB.Name)
		}

		releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
		if err != nil {
			return errors.Wrap(err, "failed to parse release date")
		}

		created, err = s.songRepo.Create(ctx, model.Song{
			GroupID:     groupDB.ID,
			Song:        song,
			Text:        songDetail.Text,
			Link:        songDetail.Link,
			ReleaseDate: releaseDate,
		})

		created.Group = groupDB.Name

		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add song")
	}

	s.logger.Info("song created",
		"id", created.ID,
		"group", created.Group,
		"song", created.Song)

	return created, nil
}

func (s *songService) Edit(ctx context.Context, song model.Song) (*model.Song, error) {
	songDB, err := s.songRepo.GetByID(ctx, song.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get song")
	}

	songDetail, err := s.musicInfoClient.GetSongInfo(ctx, song.Group, song.Song)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get song detail")
	}

	var updated *model.Song
	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		groupDB, err := s.groupRepo.GetByName(ctx, song.Group)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return errors.Wrap(err, "failed to get group by name")
		}

		if groupDB == nil {
			groupDB, err = s.groupRepo.Create(ctx, model.Group{Name: song.Group})
			if err != nil {
				return err
			}

			s.logger.Info("group created", "id", groupDB.ID, "group", groupDB.Name)
		}

		releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
		if err != nil {
			return errors.Wrap(err, "failed to parse release date")
		}

		if song.GroupID == uuid.Nil {
			song.GroupID = groupDB.ID
		}
		if song.Group == "" {
			song.Group = songDB.Group
		}
		if song.Song == "" {
			song.Song = songDB.Song
		}
		if song.Text == "" {
			song.Text = songDetail.Text
		}
		if song.Link == "" {
			song.Link = songDetail.Link
		}
		if song.ReleaseDate.IsZero() {
			song.ReleaseDate = releaseDate
		}

		song.CreatedAt = songDB.CreatedAt
		song.UpdatedAt = time.Now()

		updated, err = s.songRepo.Update(ctx, song)
		updated.Group = groupDB.Name

		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to edit song")
	}

	s.logger.Info("song updated",
		"id", updated.ID,
		"group", updated.Group,
		"song", updated.Song)

	return updated, nil
}

func (s *songService) Delete(ctx context.Context, id uuid.UUID) (*model.Song, error) {
	song, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get song")
	}

	err = s.songRepo.Delete(ctx, *song)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete song")
	}

	s.logger.Info("song deleted", "id", song.ID, "group", song.Group, "song", song.Song)

	return song, nil
}
