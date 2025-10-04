package factory

import (
	"context"
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // используется в migrate.NewWithDatabaseInstance
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mondegor/go-storage/mrmigrate/gomigrate"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrpostgres/monitoring"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-storage/mrstorage/mrprometheus"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
)

// InitPostgres - создаёт объект mrpostgres.ConnAdapter.
func InitPostgres(opts app.Options) (*mrpostgres.ConnAdapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mrlog.Info(opts.Logger, "Create and init postgres pool connection")

	cfg := opts.Cfg.Storage

	postgresOpts := mrpostgres.Options{
		Host:            cfg.Host,
		Port:            cfg.Port,
		Database:        cfg.Database,
		Username:        cfg.Username,
		Password:        cfg.Password,
		MaxPoolSize:     cfg.MaxPoolSize,
		MaxConnLifetime: cfg.MaxConnLifetime,
		MaxConnIdleTime: cfg.MaxConnIdleTime,
		ConnTimeout:     cfg.Timeout,
		QueryTracer:     monitoring.NewQueryTracer(cfg.Host, cfg.Port, cfg.Database, opts.Tracer),
	}

	conn := mrpostgres.New()

	if err := conn.Connect(ctx, postgresOpts); err != nil {
		return nil, err
	}

	connCli, err := conn.Cli()
	if err != nil {
		return nil, err
	}

	opts.Prometheus.Add(
		mrprometheus.NewDBCollector(
			"pgx",
			func() mrstorage.DBStatProvider {
				return connCli.Stat()
			},
			map[string]string{
				"db_name": cfg.Database,
			},
		),
	)

	return conn, conn.Ping(ctx)
}

// InitPostgresConnManager - создаёт объект mrpostgres.ConnManager.
func InitPostgresConnManager(conn *mrpostgres.ConnAdapter, logger mrlog.Logger) *mrpostgres.ConnManager {
	return mrpostgres.NewConnManager(conn, logger)
}

// ApplyPostgresMigrations - накатывает миграции БД opts.Cfg.Storage.
func ApplyPostgresMigrations(opts app.Options) error {
	mrlog.Info(opts.Logger, "Apply postgres migrations: "+opts.Cfg.Storage.MigrationsDir)

	connCli, err := opts.PostgresConnManager.ConnAdapter().Cli()
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(connCli)
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

	dbMigrate.Log = gomigrate.NewLoggerAdapter(opts.Logger.WithAttrs("migrator", "go-migrate"))

	if err = dbMigrate.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
