package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	authcfg "github.com/mondegor/go-components/factory/mrauth/config"
	"github.com/mondegor/go-sysmess/mrapp"
	extfilecfg "github.com/mondegor/go-sysmess/mrlib/extfile/config"
	accesscfg "github.com/mondegor/go-webcore/mraccess/config"
)

const (
	detectVersion     = "v0.0.0"
	defaultConfigPath = "./config/config.yaml"
)

// Create - создаёт, инициализирует и возвращает конфигурацию приложения.
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

	// инициализируется ENV окружение
	if err = prepareEnv(args); err != nil {
		return Config{}, fmt.Errorf("error prepare ENV: %w", err)
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
		if ver := mrapp.Version(); ver != "" {
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

		if !path.IsAbs(cfg.AccessControl.RolesDirPath) {
			cfg.AccessControl.RolesDirPath = path.Join(args.WorkDir, cfg.AccessControl.RolesDirPath)
		}
	}

	// дополнительно проверяется загруженная конфигурация
	if err = authcfg.ValidateRealms(cfg.AccessControl.Realms, cfg.AccessControl.Roles); err != nil {
		return Config{}, err
	}

	cfg.AccessControl.Realms = authcfg.DefaultValuesRealm(cfg.AccessControl.Realms, cfg.AccessControl.OperationConfirm)

	if err = accesscfg.ValidateActionGroups(cfg.AccessControl.ActionGroups, cfg.AccessControl.Privileges); err != nil {
		return Config{}, err
	}

	if err = extfilecfg.ValidateMimeTypes(cfg.Validation.MimeTypes); err != nil {
		return Config{}, err
	}

	if err = authcfg.ValidateJWT(cfg.AccessControl); err != nil {
		return Config{}, err
	}

	if len(cfg.Localization.Languages) == 0 {
		cfg.Localization.Languages = append(cfg.Localization.Languages, "ru-RU")
	}

	if cfg.Debugging.UnexpectedHttpStatus < 400 || cfg.Debugging.UnexpectedHttpStatus > 599 {
		return Config{}, fmt.Errorf("unexpected_http_status: min=400, max=599, got=%d", cfg.Debugging.UnexpectedHttpStatus)
	}

	if !cfg.Debugging.Debug {
		cfg.Debugging.UnexpectedHttpStatus = 500 // http.StatusInternalServerError
	}

	return cfg, nil
}

func prepareEnv(args Args) error {
	// загружаются ENV переменные из .env файла, если он был указан
	if args.DotEnvPath != "" {
		if err := godotenv.Overload(args.DotEnvPath); err != nil {
			return fmt.Errorf("error reading ENV file '%s': %w", args.DotEnvPath, err)
		}
	}

	// данный код, удаляет двойные кавычки у ENV значений, если они встречаются
	// это было сделано, потому как ENV, которые приходят через докер, остаются в кавычках
	for _, item := range os.Environ() {
		keyValue := strings.SplitN(item, "=", 2)

		// если не указаны ключ/значение, или значение пустое
		if len(keyValue) != 2 || len(keyValue[1]) < 2 {
			continue
		}

		// если значение не заключено в двойные кавычки
		if keyValue[1][0] != '"' || keyValue[1][len(keyValue[1])-1] != '"' {
			continue
		}

		value, err := strconv.Unquote(keyValue[1])
		if err != nil {
			return fmt.Errorf("strconv.Unquote: %w", err)
		}

		_ = os.Setenv(keyValue[0], value)
	}

	return nil
}
