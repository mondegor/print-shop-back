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
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/config"
)

// InitPostgres - создаёт объект mrpostgres.ConnAdapter.
func InitPostgres(logger mrlog.Logger, tracer mrtrace.Tracer, cfg config.Config) (*mrpostgres.ConnAdapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mrlog.Info(logger, "Create and init postgres pool connection", "host", cfg.DBHost, "hort", cfg.DBPort)

	postgresOpts := mrpostgres.Options{
		Host:            cfg.DBHost,
		Port:            cfg.DBPort,
		Database:        cfg.DBDatabase,
		Username:        cfg.DBUsername,
		Password:        cfg.DBPassword,
		MaxPoolSize:     int(cfg.DBMaxPoolSize),
		MaxConnLifetime: cfg.DBMaxConnLifetime,
		MaxConnIdleTime: cfg.DBMaxConnIdleTime,
		ConnTimeout:     cfg.DBTimeout,
		QueryTracer:     monitoring.NewQueryTracer(cfg.DBHost, cfg.DBPort, cfg.DBDatabase, tracer),
	}

	conn := mrpostgres.New()

	if err := conn.Connect(ctx, postgresOpts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}

// InitPostgresConnManager - создаёт объект mrpostgres.ConnManager.
func InitPostgresConnManager(conn *mrpostgres.ConnAdapter, logger mrlog.Logger) *mrpostgres.ConnManager {
	return mrpostgres.NewConnManager(conn, logger)
}

// ApplyPostgresMigrations - накатывает миграции БД opts.Cfg.DBMigrationsDir.
func ApplyPostgresMigrations(logger mrlog.Logger, connManager *mrpostgres.ConnManager, cfg config.Config) error {
	if cfg.DBMigrationsDir == "" {
		return nil
	}

	mrlog.Info(logger, "Apply postgres migrations", "dir", cfg.DBMigrationsDir)

	connCli, err := connManager.ConnAdapter().Cli()
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(connCli)

	defer func() {
		_ = db.Close()
	}()

	// if opts.Cfg.DBMigrationsTable is empty then will be used postgres.DefaultMigrationsTable
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: cfg.DBMigrationsTable})
	if err != nil {
		return err
	}

	dbMigrate, err := migrate.NewWithDatabaseInstance("file://"+cfg.DBMigrationsDir, cfg.DBDatabase, driver)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = dbMigrate.Close()
	}()

	dbMigrate.Log = gomigrate.NewLoggerAdapter(mrlog.WithAttrs(logger, "migrator", "go-migrate"))

	if err = dbMigrate.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
