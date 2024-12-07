package service

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math"
	"song-library-api/src/internal/model"
	"song-library-api/src/internal/repository"
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
}

func NewSongService(
	songRepo repository.SongRepository,
	groupRepo repository.GroupRepository,
	musicInfoClient *music_info_client.MusicInfoClient,
	trManager *manager.Manager) *songService {
	return &songService{
		songRepo:        songRepo,
		groupRepo:       groupRepo,
		musicInfoClient: musicInfoClient,
		trManager:       trManager,
	}
}

func (s *songService) GetSongs(ctx context.Context, filters *model.SongFilter, page, pageSize uint) (*model.PaginatedList[model.Song], error) {
	songs, err := s.songRepo.GetSongs(ctx, filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	total, err := s.songRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	totalPages := uint(math.Ceil(float64(total) / float64(pageSize)))

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

	verses := strings.Split(song.Song, "\n\n")
	total := uint(len(verses))
	totalPages := uint(math.Ceil(float64(total) / float64(pageSize)))

	start := (page - 1) * pageSize
	if start > total {
		return nil, errors.Wrap(model.ErrBadRequest, "page out of range")
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return &model.PaginatedList[string]{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Items:      verses[start:end],
	}, nil
}

func (s *songService) GetByID(ctx context.Context, id uuid.UUID) (*model.Song, error) {
	return s.songRepo.GetByID(ctx, id)
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

	return updated, nil
}

func (s *songService) Delete(ctx context.Context, id uuid.UUID) (*model.Song, error) {
	song, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get song")
	}

	return song, s.songRepo.Delete(ctx, *song)
}
