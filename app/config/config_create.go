package config

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/mondegor/go-sysmess/mrapp"
	"github.com/mondegor/go-sysmess/mrlib/extfile"

	"github.com/mondegor/print-shop-back/internal/factory/auth"
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
	if err = validateRealms(cfg.AccessControl.Realms, cfg.AccessControl.Roles); err != nil {
		return Config{}, err
	}

	cfg.AccessControl.Realms = defaultValuesRealm(cfg.AccessControl.Realms, cfg.AccessControl.OperationConfirm)

	if err = validateRoutingSections(cfg.AccessControl.RoutingSections, cfg.AccessControl.Privileges); err != nil {
		return Config{}, err
	}

	if err = validateMimeTypes(cfg.Validation.MimeTypes); err != nil {
		return Config{}, err
	}

	if err = validateJWT(cfg.AccessControl); err != nil {
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

func validateRealms(realms []auth.UserRealm, allRoles []string) error {
	uniqRealms := make(map[string]struct{}, len(realms))

	for _, realm := range realms {
		if _, ok := uniqRealms[realm.Name]; ok {
			return fmt.Errorf("duplicate realm name '%s'", realm.Name)
		}

		if realm.RegisterUserKind == "" {
			return fmt.Errorf("registerUser is empty for realm '%s'", realm.Name)
		}

		if realm.AuthToken.AccessType != "jwt" && realm.AuthToken.AccessType != "session" {
			return fmt.Errorf("invalid token type '%s' for realm '%s'", realm.AuthToken.AccessType, realm.Name)
		}

		uniqRealms[realm.Name] = struct{}{}

		if err := validateRealm(realm, allRoles); err != nil {
			return err
		}
	}

	return nil
}

func validateRealm(realm auth.UserRealm, allRoles []string) error {
	uniqKinds := make(map[string]struct{}, len(realm.UserKinds))
	hasRegisterUser := realm.RegisterUserKind == "none"

	for _, kind := range realm.UserKinds {
		if _, ok := uniqKinds[kind.Name]; ok {
			return fmt.Errorf("duplicate user kind name '%s' for realm '%s'", kind.Name, realm.Name)
		}

		uniqKinds[kind.Name] = struct{}{}

		if realm.RegisterUserKind == kind.Name {
			hasRegisterUser = true
		}

		for _, role := range kind.Roles {
			if !stringInArray(role, allRoles) {
				return fmt.Errorf("role '%s' of user kind '%s' is not found in roles for realm '%s'", role, kind.Name, realm.Name)
			}
		}
	}

	if !hasRegisterUser {
		return fmt.Errorf("realm.RegisterUserKind '%s' is not found in realm.UserKinds for realm: %s", realm.RegisterUserKind, realm.Name)
	}

	return nil
}

func defaultValuesRealm(realms []auth.UserRealm, dop auth.OperationConfirm) []auth.UserRealm {
	for i := range realms {
		rop := &realms[i].OperationConfirm

		if rop.TokenLength < 1 {
			rop.TokenLength = dop.TokenLength
		}

		if rop.CodeLength < 1 {
			rop.CodeLength = dop.CodeLength
		}

		if rop.SessionExpiry < 1 {
			rop.SessionExpiry = dop.SessionExpiry
		}

		rop.SendByEmail = defaultValuesCodeSender(rop.SendByEmail, dop.SendByEmail)
		rop.SendByPhone = defaultValuesCodeSender(rop.SendByPhone, dop.SendByPhone)
	}

	return realms
}

func defaultValuesCodeSender(cs, dcs auth.CodeSender) auth.CodeSender {
	if cs.MaxAttempts < 1 {
		cs.MaxAttempts = dcs.MaxAttempts
	}

	if cs.MaxResends < 1 {
		cs.MaxResends = dcs.MaxResends
	}

	if cs.MinResendTime < 1 {
		cs.MinResendTime = dcs.MinResendTime
	}

	return cs
}

func validateJWT(accessControl AccessControl) error {
	if !isJWTUsed(accessControl.Realms) {
		return nil
	}

	if accessControl.JWTMethod == "" {
		return errors.New("JWT method is required")
	}

	switch accessControl.JWTMethod {
	case "HS256", "HS512": // TODO: "ES256", "ES512"
	default:
		return errors.New("invalid JWT method")
	}

	if accessControl.JWTSecret == "" {
		return errors.New("JWT secret is required")
	}

	return nil
}

func isJWTUsed(realms []auth.UserRealm) bool {
	for _, realm := range realms {
		if realm.AuthToken.AccessType == "jwt" {
			return true
		}
	}

	return false
}

func validateRoutingSections(sections []RoutingSection, allPrivileges []string) error {
	uniqNames := make(map[string]struct{}, len(sections))
	uniqPaths := make(map[string]struct{}, len(sections))

	for _, section := range sections {
		if _, ok := uniqNames[section.Name]; ok {
			return fmt.Errorf("duplicate section name '%s'", section.Name)
		}

		if _, ok := uniqPaths[section.BasePath]; ok {
			return fmt.Errorf("duplicate base path '%s' for section '%s'", section.BasePath, section.Name)
		}

		uniqNames[section.Name] = struct{}{}
		uniqPaths[section.BasePath] = struct{}{}

		if section.Privilege != "public" {
			if !stringInArray(section.Privilege, allPrivileges) {
				return fmt.Errorf("'%s' is not found in privileges for section '%s'", section.Privilege, section.Name)
			}
		}
	}

	return nil
}

func validateMimeTypes(mimeTypes []extfile.MimeType) error {
	uniqExtensions := make(map[string]struct{}, len(mimeTypes))

	for _, mimeType := range mimeTypes {
		if _, ok := uniqExtensions[mimeType.Extension]; ok {
			return fmt.Errorf("duplicate mimeType extension '%s'", mimeType.Extension)
		}

		uniqExtensions[mimeType.Extension] = struct{}{}
	}

	return nil
}

func stringInArray(needle string, haystack []string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}
