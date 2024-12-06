package app

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"log"
	"song-library-api/src/internal/config"
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
	m, err := migrate.New("file://./src/cmd/migrate/migrations/postgresql", a.provider.Config().PostgresConn)
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
	group := a.httpServer.Group("/api")

	a.httpServer.Use(middleware.Recover())
	a.httpServer.Use(middleware.Logger())

	return nil
}

func (a *App) runHttpServer() error {
	err := a.httpServer.Start(a.provider.Config().ServerAddress)
	if err != nil {
		return errors.Wrap(err, "failed to start http server")
	}
	return nil
}
