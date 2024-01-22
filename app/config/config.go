package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	appName    = "Print Shop"
	appVersion = "v0.8.0"
)

type (
	Config struct {
		AppName         string
		AppVersion      string
		AppInfo         string
		AppPath         string
		AppStartedAt    time.Time
		ConfigPath      string
		Debugging       `yaml:"debugging"`
		Log             `yaml:"logger"`
		Server          `yaml:"server"`
		Listen          `yaml:"listen"`
		Storage         `yaml:"storage"`
		Redis           `yaml:"redis"`
		FileSystem      `yaml:"file_system"`
		FileProviders   `yaml:"file_providers"`
		Cors            `yaml:"cors"`
		Translation     `yaml:"translation"`
		AppSections     `yaml:"app_sections"`
		AccessControl   `yaml:"access_control"`
		ModulesSettings `yaml:"modules_settings"`
	}

	Debugging struct {
		Debug       bool `yaml:"debug" env:"APPX_DEBUG"`
		ErrorCaller `yaml:"caller"`
	}

	ErrorCaller struct {
		Deep         int    `yaml:"deep" env:"APPX_ERR_CALLER_DEEP"`
		UseShortPath bool   `yaml:"use_short_path" env:"APPX_ERR_CALLER_USE_SHORT_PATH"`
		RootPath     string `yaml:"root_path"`
	}

	Log struct {
		Prefix    string `yaml:"prefix" env:"APPX_LOG_PREFIX"`
		Level     string `yaml:"level" env:"APPX_LOG_LEVEL"`
		LogCaller `yaml:"caller"`
	}

	LogCaller struct {
		Deep         int    `yaml:"deep" env:"APPX_LOG_CALLER_DEEP"`
		UseShortPath bool   `yaml:"use_short_path" env:"APPX_LOG_CALLER_USE_SHORT_PATH"`
		RootPath     string `yaml:"root_path"`
	}

	Server struct {
		ReadTimeout     time.Duration `yaml:"read_timeout" env:"APPX_SERVER_READ_TIMEOUT"`
		WriteTimeout    time.Duration `yaml:"write_timeout" env:"APPX_SERVER_WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"APPX_SERVER_SHUTDOWN_TIMEOUT"`
	}

	Listen struct {
		Type     string `yaml:"type" env:"APPX_SERVICE_LISTEN_TYPE"`
		SockName string `yaml:"sock_name" env:"APPX_SERVICE_LISTEN_SOCK"`
		BindIP   string `yaml:"bind_ip" env:"APPX_SERVICE_BIND"`
		Port     string `yaml:"port" env:"APPX_SERVICE_PORT"`
	}

	Storage struct {
		Host        string        `yaml:"host" env:"APPX_DB_HOST"`
		Port        string        `yaml:"port" env:"APPX_DB_PORT"`
		Username    string        `yaml:"username" env:"APPX_DB_USER"`
		Password    string        `yaml:"password" env:"APPX_DB_PASSWORD"`
		Database    string        `yaml:"database" env:"APPX_DB_NAME"`
		MaxPoolSize int           `yaml:"max_pool_size" env:"APPX_DB_MAX_POOL_SIZE"`
		Timeout     time.Duration `yaml:"timeout"`
	}

	Redis struct {
		Host     string        `yaml:"host" env:"APPX_REDIS_HOST"`
		Port     string        `yaml:"port" env:"APPX_REDIS_PORT"`
		Password string        `yaml:"password" env:"APPX_REDIS_PASSWORD"`
		Timeout  time.Duration `yaml:"timeout"`
	}

	FileSystem struct {
		DirMode    uint32 `yaml:"dir_mode" env:"APPX_FILESYSTEM_DIR_MODE"`
		CreateDirs bool   `yaml:"create_dirs" env:"APPX_FILESYSTEM_CREATE_DIRS"`
	}

	FileProviders struct {
		ImageStorage struct {
			Name    string `yaml:"name"`
			RootDir string `yaml:"root_dir" env:"APPX_IMAGESTORAGE_ROOT_DIR"`
		} `yaml:"image_storage"`
	}

	Cors struct {
		AllowedOrigins   []string `yaml:"allowed_origins" env:"APPX_CORS_ALLOWED_ORIGINS"` // items by "," separated
		AllowedMethods   []string `yaml:"allowed_methods"`
		AllowedHeaders   []string `yaml:"allowed_headers"`
		ExposedHeaders   []string `yaml:"exposed_headers"`
		AllowCredentials bool     `yaml:"allow_credentials"`
	}

	Translation struct {
		DirPath      string   `yaml:"dir_path"`
		LangCodes    []string `yaml:"lang_codes" env:"APPX_TRANSLATION_LANGS"` // items by "," separated
		Dictionaries struct {
			DirPath string   `yaml:"dir_path"`
			List    []string `yaml:"list"`
		} `yaml:"dictionaries"`
	}

	AppSections struct {
		AdminAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_ADMIN_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_ADMIN_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"admin_api"`
		ProviderAccountAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_PR_ACCOUNT_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_PR_ACCOUNT_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"provider_account_api"`
		PublicAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_PUBLIC_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_PUBLIC_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"public_api"`
	}

	AccessControl struct {
		Roles       `yaml:"roles"`
		Privileges  []string `yaml:"privileges"`
		Permissions []string `yaml:"permissions"`
	}

	Roles struct {
		DirPath  string   `yaml:"dir_path"`
		FileType string   `yaml:"file_type"`
		List     []string `yaml:"list"`
	}

	ModulesSettings struct {
		ProviderAccount struct {
			CompanyPageLogo struct {
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage
			} `yaml:"company_page_logo"`
		} `yaml:"provider_account"`
		FileStation struct {
			ImageProxy struct {
				Host         string `yaml:"host" env:"APPX_IMAGE_HOST"`
				BaseURL      string `yaml:"base_url"`
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage
			} `yaml:"image_proxy"`
		} `yaml:"file_station"`
	}
)

func New(filePath string) (*Config, error) {
	cfg := &Config{
		AppName:    appName,
		AppVersion: appVersion,
		ConfigPath: filePath,
	}

	if err := cleanenv.ReadConfig(filePath, cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file '%s': %w", filePath, err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("error reading ENV from config file '%s': %w", filePath, err)
	}

	cfg.AppPath = os.Args[0]
	cfg.AppStartedAt = time.Now().UTC()

	return cfg, nil
}
