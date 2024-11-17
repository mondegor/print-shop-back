package factory

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // используется в migrate.NewWithDatabaseInstance
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mondegor/go-storage/mrmigrate/gomigrate"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrpostgres/logger"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-storage/mrstorage/mrprometheus"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
)

// NewPostgres - создаёт объект mrpostgres.ConnAdapter.
func NewPostgres(ctx context.Context, opts app.Options) (*mrpostgres.ConnAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init postgres connection")

	cfg := opts.Cfg
	postgresOpts := mrpostgres.Options{
		Host:            cfg.Storage.Host,
		Port:            cfg.Storage.Port,
		Database:        cfg.Storage.Database,
		Username:        cfg.Storage.Username,
		Password:        cfg.Storage.Password,
		MaxPoolSize:     cfg.Storage.MaxPoolSize,
		MaxConnLifetime: cfg.Storage.MaxConnLifetime,
		MaxConnIdleTime: cfg.Storage.MaxConnIdleTime,
		ConnAttempts:    1,
		ConnTimeout:     cfg.Storage.Timeout,
		QueryTracer:     logger.NewQueryTracer(cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database),
	}

	conn := mrpostgres.New()

	if err := conn.Connect(ctx, postgresOpts); err != nil {
		return nil, err
	}

	opts.Prometheus.MustRegister(
		mrprometheus.NewDBCollector(
			"pgx",
			func() mrstorage.DBStatProvider {
				return conn.Cli().Stat()
			},
			map[string]string{
				"db_name": cfg.Storage.Database,
			},
		),
	)

	return conn, conn.Ping(ctx)
}

// NewPostgresConnManager - создаёт объект mrpostgres.ConnManager.
func NewPostgresConnManager(_ context.Context, conn *mrpostgres.ConnAdapter) *mrpostgres.ConnManager {
	return mrpostgres.NewConnManager(conn)
}

// ApplyPostgresMigrations - накатывает миграции БД opts.Cfg.Storage.
func ApplyPostgresMigrations(ctx context.Context, opts app.Options) error {
	mrlog.Ctx(ctx).Info().Msgf("Apply postgres migrations: %s", opts.Cfg.Storage.MigrationsDir)

	db := stdlib.OpenDBFromPool(opts.PostgresConnManager.ConnAdapter().Cli())
	defer db.Close()

	// if opts.Cfg.Storage.MigrationsTable is empty then will be used postgres.DefaultMigrationsTable
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: opts.Cfg.Storage.MigrationsTable})
	if err != nil {
		return err
	}

	dbMigrate, err := migrate.NewWithDatabaseInstance("file://"+opts.Cfg.Storage.MigrationsDir, opts.Cfg.Storage.Database, driver)
	if err != nil {
		return err
	}
	defer dbMigrate.Close()

	dbMigrate.Log = gomigrate.NewLoggerAdapter(mrlog.Ctx(ctx).With().Str("migrator", "go-migrate").Logger())

	if err = dbMigrate.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
