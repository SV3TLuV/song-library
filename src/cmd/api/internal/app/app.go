package app

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"song-library-api/src/cmd/api/internal/config"
	middleware2 "song-library-api/src/cmd/api/internal/server/http/middleware"
	"song-library-api/src/cmd/api/internal/server/http/route"
	v1 "song-library-api/src/cmd/api/internal/server/http/v1"
	"song-library-api/src/cmd/api/internal/server/http/validator"
)

type App struct {
	provider   *serviceProvider
	httpServer *echo.Echo
}

func New() (*App, error) {
	a := &App{}
	ctx := context.Background()
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	return a.runHttpServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHttpServer,
		a.applyMigration,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) applyMigration(_ context.Context) error {
	m, err := migrate.New("file://./src/cmd/api/migrations/postgresql", a.provider.Config().PostgresConn)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil {
		return err
	}

	log.Printf("Applied migration: %d. Dirty: %t\n", version, dirty)
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	return config.Load()
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.provider = NewServiceProvider()
	return nil
}

func (a *App) initHttpServer(_ context.Context) error {
	a.httpServer = echo.New()

	a.httpServer.Use(middleware.Recover())
	a.httpServer.Use(middleware.Logger())
	a.httpServer.Use(middleware2.ErrorHandlerMiddleware)

	a.httpServer.Validator = validator.NewRequestValidator()

	group := a.httpServer.Group("")
	route.InitSongRoutes(group, v1.NewSongController(a.provider.SongService()))

	a.httpServer.GET("/swagger/*", echoSwagger.WrapHandler)

	return nil
}

func (a *App) runHttpServer() error {
	err := a.httpServer.Start(a.provider.Config().ServerAddress)
	if err != nil {
		return errors.Wrap(err, "failed to start http server")
	}
	return nil
}
