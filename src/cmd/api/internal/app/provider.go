package app

import (
	"context"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"log"
	"log/slog"
	"os"
	"song-library-api/src/cmd/api/internal/config"
	"song-library-api/src/cmd/api/internal/db/postgres"
	"song-library-api/src/cmd/api/internal/repository"
	"song-library-api/src/cmd/api/internal/service"
	"song-library-api/src/pkg/music_info_client"
)

type serviceProvider struct {
	config *config.Config

	logger *slog.Logger

	postgres  *pgxpool.Pool
	trManager *manager.Manager

	musicInfoClient *music_info_client.MusicInfoClient

	songRepo    repository.SongRepository
	songService service.SongService

	groupRepo    repository.GroupRepository
	groupService service.GroupService
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (p *serviceProvider) Config() *config.Config {
	if p.config == nil {
		cfg, err := config.FromEnv()
		if err != nil {
			log.Fatal(errors.Wrap(err, "init config"))
		}
		p.config = cfg
	}
	return p.config
}

func (p *serviceProvider) Logger() *slog.Logger {
	if p.logger == nil {
		p.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
	return p.logger
}

func (p *serviceProvider) Postgres() *pgxpool.Pool {
	if p.postgres == nil {
		db, err := postgres.NewDB(context.Background(), p.Config().PostgresConn)
		if err != nil {
			log.Fatal(errors.Wrap(err, "init postgresql pool"))
		}
		p.postgres = db
	}
	return p.postgres
}

func (p *serviceProvider) TransactionManager() *manager.Manager {
	if p.trManager == nil {
		p.trManager = manager.Must(trmpgx.NewDefaultFactory(p.postgres))
	}
	return p.trManager
}

func (p *serviceProvider) MusicInfoClient() *music_info_client.MusicInfoClient {
	if p.musicInfoClient == nil {
		p.musicInfoClient = music_info_client.NewMusicInfoClient(p.Config().MusicInfoServiceURL)
	}
	return p.musicInfoClient
}

func (p *serviceProvider) SongRepo() repository.SongRepository {
	if p.songRepo == nil {
		p.songRepo = repository.NewSongRepository(p.Postgres(),
			trmpgx.DefaultCtxGetter,
			p.TransactionManager())
	}
	return p.songRepo
}

func (p *serviceProvider) GroupRepo() repository.GroupRepository {
	if p.groupRepo == nil {
		p.groupRepo = repository.NewGroupRepository(p.Postgres(),
			trmpgx.DefaultCtxGetter,
			p.TransactionManager())
	}
	return p.groupRepo
}

func (p *serviceProvider) SongService() service.SongService {
	if p.songService == nil {
		p.songService = service.NewSongService(
			p.SongRepo(),
			p.GroupRepo(),
			p.MusicInfoClient(),
			p.TransactionManager(),
			p.Logger())
	}
	return p.songService
}

func (p *serviceProvider) GroupService() service.GroupService {
	if p.groupService == nil {
		p.groupService = service.NewGroupService(
			p.GroupRepo(),
			p.TransactionManager(),
			p.Logger())
	}
	return p.groupService
}
