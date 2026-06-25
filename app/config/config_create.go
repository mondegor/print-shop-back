package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	authcfg "github.com/mondegor/go-components/wire/mrauth/config"
	accesscfg "github.com/mondegor/go-sysmess/mraccess/config"
	"github.com/mondegor/go-sysmess/mrapp"
	extfilecfg "github.com/mondegor/go-sysmess/util/mime/config"
)

const (
	environmentName    = "APPX_ENV"
	defaultEnvironment = "prod"
	detectVersion      = "v0.0.0"
	relativeConfigDir  = "./config/"
)

var regexpEnvironment = regexp.MustCompile(`^[a-z][a-z0-9]+$`)

// Create - создаёт, инициализирует и возвращает конфигурацию приложения.
func Create(args CmdArgs, stdout io.Writer) (cfg Config, err error) {
	if stdout == nil {
		return Config{}, errors.New("stdout is required")
	}

	if cfg, err = createConfig(args); err != nil {
		return Config{}, fmt.Errorf("error readConfig: %w", err)
	}

	cfg.Stdout = stdout

	cfg.AccessControl.Realms = authcfg.CorrectValuesRealm(
		cfg.AccessControl.Realms,
		cfg.AccessControl.DefaultOperationConfirm,
		cfg.AccessControl.OverrideAuthToken,
	)

	// дополнительно проверяется загруженная конфигурация
	if err = authcfg.ValidateRealms(cfg.AccessControl.Realms, cfg.AccessControl.Roles); err != nil {
		return Config{}, err
	}

	if err = accesscfg.ValidateActionGroups(cfg.AccessControl.ActionGroups); err != nil {
		return Config{}, err
	}

	if err = extfilecfg.ValidateMimeTypes(cfg.AllowedMimeTypes); err != nil {
		return Config{}, err
	}

	if authcfg.IsJWTUsed(cfg.AccessControl.Realms) {
		if cfg.JWT, err = authcfg.InitJWT(cfg.JWT); err != nil {
			return Config{}, err
		}
	}

	if cfg.UnexpectedErrorHttpStatus < 400 || cfg.UnexpectedErrorHttpStatus > 599 {
		return Config{}, fmt.Errorf("unexpected_error_http_status: min=400, max=599, got=%d", cfg.UnexpectedErrorHttpStatus)
	}

	if cfg.AppVersion == detectVersion {
		if ver := mrapp.Version(); ver != "" {
			cfg.AppVersion = ver
		}
	}

	if len(cfg.AppLanguages) == 0 {
		cfg.AppLanguages = append(cfg.AppLanguages, "ru-RU")
	}

	if cfg.WorkDir != "" {
		if cfg.DBMigrationsDir != "" && !path.IsAbs(cfg.DBMigrationsDir) {
			cfg.DBMigrationsDir = path.Join(cfg.WorkDir, cfg.DBMigrationsDir)
		}

		if !path.IsAbs(cfg.FileProviders.ImageStorage2RootDir) {
			cfg.FileProviders.ImageStorage2RootDir = path.Join(cfg.WorkDir, cfg.FileProviders.ImageStorage2RootDir)
		}

		if !path.IsAbs(cfg.AccessControl.RolesDirPath) {
			cfg.AccessControl.RolesDirPath = path.Join(args.WorkDir, cfg.AccessControl.RolesDirPath)
		}
	}

	return cfg, nil
}

func createConfig(args CmdArgs) (cfg Config, err error) {
	cfg.StartedAt = time.Now().UTC()

	// загружаются ENV переменные из .env файла, если он был явно указан
	if args.DotEnvPath != "" {
		if args.WorkDir != "" && !path.IsAbs(args.DotEnvPath) {
			args.DotEnvPath = path.Join(args.WorkDir, args.DotEnvPath)
		}

		if err := godotenv.Overload(args.DotEnvPath); err != nil {
			return Config{}, fmt.Errorf("error reading ENV file '%s': %w", args.DotEnvPath, err)
		}

		cfg.DotEnvPath = args.DotEnvPath
	}

	if err := unquoteEnvs(); err != nil {
		return Config{}, fmt.Errorf("error unquote Envs: %w", err)
	}

	if args.Environment == "" {
		args.Environment = os.Getenv(environmentName)

		if args.Environment == "" {
			args.Environment = defaultEnvironment
		}
	}

	if !regexpEnvironment.MatchString(args.Environment) {
		return Config{}, errors.New("args.Environment must match " + regexpEnvironment.String())
	}

	cfg.Environment = args.Environment
	configDir := relativeConfigDir

	if args.WorkDir != "" {
		configDir = path.Join(args.WorkDir, configDir) + "/"
		cfg.WorkDir = args.WorkDir
	}

	// загружается базовая конфигурация
	if err = parseYAML(configDir+"config.yaml", &cfg, true); err != nil {
		return Config{}, err
	}

	// основная конфигурация уточняется, если задана конфигурация для указанного окружения
	if err = parseYAML(configDir+"config_"+cfg.Environment+".yaml", &cfg, false); err != nil {
		return Config{}, err
	}

	// основная конфигурация уточняется переменными из ENV окружения
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("error reading ENV: %w", err)
	}

	if args.LogLevel != "" {
		cfg.LogLevel = args.LogLevel
	}

	return cfg, nil
}

func parseYAML(configPath string, cfg *Config, required bool) error {
	fp, err := os.Open(configPath) //nolint:gosec
	if err != nil {
		if !required && errors.Is(err, os.ErrNotExist) {
			cfg.ConfigPaths = append(cfg.ConfigPaths, configPath+" [SKIPPED]")

			return nil
		}

		return fmt.Errorf("open config file '%s': %w", configPath, err)
	}

	defer func() {
		_ = fp.Close()
	}()

	if err = cleanenv.ParseYAML(fp, cfg); err != nil {
		return fmt.Errorf("error parsing config file '%s': %w", configPath, err)
	}

	cfg.ConfigPaths = append(cfg.ConfigPaths, configPath)

	return nil
}

// unquoteEnvs - удаляет двойные кавычки у ENV значений, если они встречаются
// (ENV приходящие через docker, могут оставаться в кавычках).
func unquoteEnvs() error {
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
