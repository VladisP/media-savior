package db

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	// db driver for applying migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// source driver for applying migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/VladisP/media-savior/internal/common/log"
	"github.com/VladisP/media-savior/internal/core/config"
)

type InputParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *config.Config
	Logger    log.ZapLogger
}

func NewDB(p InputParams) (*sqlx.DB, error) {
	dbConfig := p.Config.DB
	logger := p.Logger.NestedLogger("db")

	connString, err := buildConnConfig(dbConfig, logger)
	if err != nil {
		return nil, errors.Wrap(err, "build conn config error")
	}

	db, err := sqlx.Open("pgx", connString)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create DB abstraction")
	}
	db.SetMaxOpenConns(dbConfig.MaxOpenConnections)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err = db.Ping(); err != nil {
				logger.Error(ctx, "Failed to establish db connection", zap.Error(err))
				return err
			}
			logger.Info(ctx, "DB connection established successfully")

			if err = applyMigrations(dbConfig, logger); err != nil {
				logger.Error(ctx, "Failed to apply migrations", zap.Error(err))
				return err
			}
			logger.Info(ctx, "Migrations applied successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info(ctx, "Closing DB")
			return db.Close()
		},
	})

	return db, nil
}

func buildConnConfig(dbConfig config.DB, logger log.ZapLogger) (string, error) {
	connConfig, err := pgx.ParseConfig(buildDatabaseURL(dbConfig))
	if err != nil {
		return "", errors.Wrap(err, "parse conn config from DSN error")
	}

	connConfig.Logger = zapadapter.NewLogger(logger.Zap())
	return stdlib.RegisterConnConfig(connConfig), nil
}

func buildDatabaseURL(dbConfig config.DB) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode,
	)
}

func applyMigrations(dbConfig config.DB, logger log.ZapLogger) error {
	m, err := migrate.New(dbConfig.MigrationsDir, buildDatabaseURL(dbConfig))
	if err != nil {
		return errors.Wrap(err, "create migrate instance error")
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			logger.ErrorWithoutContext("Migrate close: source error", zap.Error(sourceErr))
		}
		if dbErr != nil {
			logger.ErrorWithoutContext("Migrate close: db error", zap.Error(dbErr))
		}
	}()

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "migrate up error")
	}
	return nil
}
