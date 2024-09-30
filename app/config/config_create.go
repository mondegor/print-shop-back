package config

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
)

const (
	detectVersion     = "v0.0.0"
	defaultConfigPath = "./config/config.yaml"
)

// Create - создаёт, инициализирует и возвращает конфигурацию приложения.
//
//nolint:nestif
func Create(args Args) (cfg Config, err error) {
	cfg.App.StartedAt = time.Now().UTC()

	if args.Stdout == nil {
		return Config{}, errors.New("args.Stdout is required")
	}

	if args.ConfigPath == "" {
		args.ConfigPath = defaultConfigPath
	}

	if args.WorkDir != "" {
		if !path.IsAbs(args.ConfigPath) {
			args.ConfigPath = path.Join(args.WorkDir, args.ConfigPath)
		}

		if args.DotEnvPath != "" && !path.IsAbs(args.DotEnvPath) {
			args.DotEnvPath = path.Join(args.WorkDir, args.DotEnvPath)
		}
	}

	// сначала загружается базовая конфигурация из config.yaml
	if err = cleanenv.ReadConfig(args.ConfigPath, &cfg); err != nil {
		return Config{}, fmt.Errorf("error parsing config file '%s': %w", args.ConfigPath, err)
	}

	// загружаются ENV переменные из .env файла, если он был указан
	if args.DotEnvPath != "" {
		if err = godotenv.Load(args.DotEnvPath); err != nil {
			return Config{}, fmt.Errorf("error reading ENV file '%s': %w", args.DotEnvPath, err)
		}
	}

	// уточняется конфигурация переменными из ENV окружения
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("error reading ENV from config file '%s': %w", args.ConfigPath, err)
	}

	// уточняется конфигурация переменными из внешнего окружения (переданные из командной строки, тестовой среды)
	cfg.App.WorkDir = args.WorkDir
	cfg.App.ConfigPath = args.ConfigPath
	cfg.App.DotEnvPath = args.DotEnvPath
	cfg.Os.Stdout = args.Stdout

	if cfg.App.Version == detectVersion {
		if ver := mrinit.Version(); ver != "" {
			cfg.App.Version = ver
		}
	}

	if args.Environment != "" {
		cfg.App.Environment = args.Environment
	}

	if args.LogLevel != "" {
		cfg.Log.Level = args.LogLevel
	}

	if args.WorkDir != "" {
		if cfg.Storage.MigrationsDir != "" && !path.IsAbs(cfg.Storage.MigrationsDir) {
			cfg.Storage.MigrationsDir = path.Join(args.WorkDir, cfg.Storage.MigrationsDir)
		}

		if !path.IsAbs(cfg.FileProviders.ImageStorage.RootDir) {
			cfg.FileProviders.ImageStorage.RootDir = path.Join(args.WorkDir, cfg.FileProviders.ImageStorage.RootDir)
		}

		if !path.IsAbs(cfg.Translation.DirPath) {
			cfg.Translation.DirPath = path.Join(args.WorkDir, cfg.Translation.DirPath)
		}

		if !path.IsAbs(cfg.Translation.Dictionaries.DirPath) {
			cfg.Translation.Dictionaries.DirPath = path.Join(args.WorkDir, cfg.Translation.Dictionaries.DirPath)
		}

		if !path.IsAbs(cfg.AccessControl.Roles.DirPath) {
			cfg.AccessControl.Roles.DirPath = path.Join(args.WorkDir, cfg.AccessControl.Roles.DirPath)
		}
	}

	// проверяется конфигурация на валидность
	if cfg.Debugging.UnexpectedHttpStatus < 400 || cfg.Debugging.UnexpectedHttpStatus > 599 {
		return Config{}, fmt.Errorf("unexpected_http_status: min=400, max=599, got=%d", cfg.Debugging.UnexpectedHttpStatus)
	}

	return cfg, nil
}
