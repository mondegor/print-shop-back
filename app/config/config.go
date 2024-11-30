package config

import (
	"io"
	"time"

	"github.com/mondegor/go-webcore/mrlib"
)

type (
	// Args - разобранные аргументы, которые передаются из командной строки.
	// Эти аргументы более приоритетны аналогичным, которые определены в конфигурации или переданы через .env файл.
	Args struct {
		WorkDir     string // путь к рабочей директории приложения
		ConfigPath  string // путь к файлу конфигурации приложения
		DotEnvPath  string // путь к .env файлу (переменные из этого файла более приоритетны переменных из ConfigPath)
		Environment string // внешнее окружение: local, dev, test, prod
		LogLevel    string // уровень логирования: info, warn, error, fatal, debug, trace
		Stdout      io.Writer
	}

	// Config - comment struct.
	Config struct {
		Os
		App             `yaml:"app"`
		Debugging       `yaml:"debugging"`
		Log             `yaml:"logger"`
		Sentry          `yaml:"sentry"`
		Servers         `yaml:"servers"`
		Storage         `yaml:"storage"`
		Redis           `yaml:"redis"`
		FileSystem      `yaml:"file_system"`
		FileProviders   `yaml:"file_providers"`
		Cors            `yaml:"cors"`
		Translation     `yaml:"translation"`
		Senders         `yaml:"senders"`
		AppSections     `yaml:"app_sections"`
		AccessControl   `yaml:"access_control"`
		ModulesSettings `yaml:"modules_settings"`
		Validation      `yaml:"validation"`
		TaskSchedule    `yaml:"task_schedule"`
	}

	// Os - comment struct.
	Os struct {
		Stdout io.Writer
	}

	// App - comment struct.
	App struct {
		Name        string `yaml:"name" env:"APPX_NAME"`
		Version     string `yaml:"version" env:"APPX_VER"`
		Environment string `yaml:"environment" env:"APPX_ENV"`
		WorkDir     string
		ConfigPath  string
		DotEnvPath  string
		StartedAt   time.Time
	}

	// Debugging - comment struct.
	Debugging struct {
		Debug                bool `yaml:"debug" env:"APPX_DEBUG"`
		UnexpectedHttpStatus int  `yaml:"unexpected_http_status"`
		ErrorCaller          `yaml:"error_caller"`
	}

	// ErrorCaller - comment struct.
	ErrorCaller struct {
		Enable       bool     `yaml:"enable" env:"APPX_ERR_CALLER_ENABLE"`
		Depth        uint8    `yaml:"depth" env:"APPX_ERR_CALLER_DEPTH"`
		ShowFuncName bool     `yaml:"show_func_name"`
		UpperBounds  []string `yaml:"upper_bounds"`
	}

	// Log - comment struct.
	Log struct {
		Level           string `yaml:"level" env:"APPX_LOG_LEVEL"`
		TimestampFormat string `yaml:"timestamp_format" env:"APPX_LOG_TIMESTAMP"`
		JsonFormat      bool   `yaml:"json_format" env:"APPX_LOG_JSON"`
		ConsoleColor    bool   `yaml:"console_color" env:"APPX_LOG_COLOR"`
	}

	// Sentry - comment struct.
	Sentry struct {
		DSN              string        `yaml:"dsn" env:"APPX_SENTRY_DSN"`
		TracesSampleRate float64       `yaml:"traces_sample_rate" env:"APPX_SENTRY_TRACES_SAMPLE_RATE"`
		FlushTimeout     time.Duration `yaml:"flush_timeout"`
	}

	// Servers - comment struct.
	Servers struct {
		// RestServer - comment struct.
		RestServer struct {
			ReadTimeout     time.Duration `yaml:"read_timeout" env:"APPX_SERVER_READ_TIMEOUT"`
			WriteTimeout    time.Duration `yaml:"write_timeout" env:"APPX_SERVER_WRITE_TIMEOUT"`
			ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"APPX_SERVER_SHUTDOWN_TIMEOUT"`
			Listen          struct {
				BindIP string `yaml:"bind_ip" env:"APPX_SERVER_LISTEN_BIND"`
				Port   string `yaml:"port" env:"APPX_SERVER_LISTEN_PORT"`
			} `yaml:"listen"`
		} `yaml:"rest_server"`

		// InternalServer - comment struct.
		InternalServer struct {
			ReadTimeout     time.Duration `yaml:"read_timeout" env:"APPX_INTERNAL_SERVER_READ_TIMEOUT"`
			WriteTimeout    time.Duration `yaml:"write_timeout" env:"APPX_INTERNAL_SERVER_WRITE_TIMEOUT"`
			ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"APPX_INTERNAL_SERVER_SHUTDOWN_TIMEOUT"`
			Listen          struct {
				BindIP string `yaml:"bind_ip" env:"APPX_INTERNAL_SERVER_LISTEN_BIND"`
				Port   string `yaml:"port" env:"APPX_INTERNAL_SERVER_LISTEN_PORT"`
			} `yaml:"listen"`
		} `yaml:"internal_server"`
	}

	// Storage - comment struct.
	Storage struct {
		Host            string        `yaml:"host" env:"APPX_DB_HOST"`
		Port            string        `yaml:"port" env:"APPX_DB_PORT"`
		Username        string        `yaml:"username" env:"APPX_DB_USER"`
		Password        string        `yaml:"password" env:"APPX_DB_PASSWORD"`
		Database        string        `yaml:"database" env:"APPX_DB_NAME"`
		MigrationsDir   string        `yaml:"migrations_dir"`
		MigrationsTable string        `yaml:"migrations_table" env:"APPX_DB_MIGRATIONS_TABLE"`
		MaxPoolSize     int           `yaml:"max_pool_size" env:"APPX_DB_MAX_POOL_SIZE"`
		MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" env:"APPX_DB_MAX_CONN_LIFETIME"`
		MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time" env:"APPX_DB_MAX_CONN_IDLE_TIME"`
		Timeout         time.Duration `yaml:"timeout"`
	}

	// Redis - comment struct.
	Redis struct {
		Host         string        `yaml:"host" env:"APPX_REDIS_HOST"`
		Port         string        `yaml:"port" env:"APPX_REDIS_PORT"`
		Password     string        `yaml:"password" env:"APPX_REDIS_PASSWORD"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env:"APPX_REDIS_READ_TIMEOUT"`
		WriteTimeout time.Duration `yaml:"write_timeout" env:"APPX_REDIS_WRITE_TIMEOUT"`
	}

	// FileSystem - comment struct.
	FileSystem struct {
		DirMode    uint32 `yaml:"dir_mode" env:"APPX_FILESYSTEM_DIR_MODE"`
		CreateDirs bool   `yaml:"create_dirs" env:"APPX_FILESYSTEM_CREATE_DIRS"`
	}

	// FileProviders - comment struct.
	FileProviders struct {
		// ImageStorage - comment struct.
		ImageStorage struct {
			Name    string `yaml:"name"`
			RootDir string `yaml:"root_dir" env:"APPX_IMAGESTORAGE_ROOT_DIR"`
		} `yaml:"image_storage"`
	}

	// Cors - comment struct.
	Cors struct {
		AllowedOrigins   []string `yaml:"allowed_origins" env:"APPX_CORS_ALLOWED_ORIGINS"` // items by "," separated
		AllowedMethods   []string `yaml:"allowed_methods"`
		AllowedHeaders   []string `yaml:"allowed_headers"`
		ExposedHeaders   []string `yaml:"exposed_headers"`
		AllowCredentials bool     `yaml:"allow_credentials"`
	}

	// Translation - comment struct.
	Translation struct {
		DirPath   string   `yaml:"dir_path" env:"APPX_TRANSLATION_DIR_PATH"`
		LangCodes []string `yaml:"lang_codes" env:"APPX_TRANSLATION_LANGS"` // items by "," separated
		// Dictionaries - comment struct.
		Dictionaries struct {
			DirPath string   `yaml:"dir_path" env:"APPX_TRANSLATION_DICTIONARIES_DIR_PATH"`
			List    []string `yaml:"list"`
		} `yaml:"dictionaries"`
	}

	// Senders - comment struct.
	Senders struct {
		Mail        SenderMail        `yaml:"mail"`
		TelegramBot SenderTelegramBot `yaml:"telegram_bot"`
	}

	SenderMail struct {
		DefaultFrom  string `yaml:"default_from" env:"APPX_SENDER_MAIL_DEFAULT_FROM"`    // e-mail от которого отправляются письма [по умолчанию]
		SmtpHost     string `yaml:"smtp_host" env:"APPX_SENDER_MAIL_SMTP_HOST"`          // SMTP сервер для отправки почты (домен или IP)
		SmtpPort     string `yaml:"smtp_port" env:"APPX_SENDER_MAIL_SMTP_PORT"`          // порт SMTP сервера для отправки почты
		SmtpUserName string `yaml:"smtp_user_name" env:"APPX_SENDER_MAIL_SMTP_USERNAME"` // адрес почтового ящика на SMTP сервере
		SmtpPassword string `yaml:"smtp_password" env:"APPX_SENDER_MAIL_SMTP_PASSWORD"`  // пароль почтового ящика на SMTP сервере
	}

	// SenderTelegramBot - comment struct.
	SenderTelegramBot struct {
		Name  string `yaml:"name" env:"APPX_SENDER_TELEGRAMBOT_NAME"`
		Token string `yaml:"token" env:"APPX_SENDER_TELEGRAMBOT_TOKEN"`
	}

	// AppSections - comment struct.
	AppSections struct {
		// AdminAPI - comment struct.
		AdminAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_ADMIN_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_ADMIN_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"admin_api"`
		// ProvidersAPI - comment struct.
		ProvidersAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_PROVIDERS_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_PROVIDERS_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"providers_api"`
		// PublicAPI - comment struct.
		PublicAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_PUBLIC_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_PUBLIC_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"public_api"`
	}

	// AccessControl - comment struct.
	AccessControl struct {
		Roles       `yaml:"roles"`
		Privileges  []string `yaml:"privileges"`
		Permissions []string `yaml:"permissions"`
	}

	// Roles - comment struct.
	Roles struct {
		DirPath  string   `yaml:"dir_path" env:"APPX_ROLES_DIR_PATH"`
		FileType string   `yaml:"file_type"`
		List     []string `yaml:"list"`
	}

	// ModulesSettings - comment struct.
	ModulesSettings struct {
		// General - comment struct.
		General struct {
			PageSizeMax     uint64 `yaml:"page_size_max"`
			PageSizeDefault uint64 `yaml:"page_size_default"`
		} `yaml:"general"`
		// ProviderAccount - comment struct.
		ProviderAccount struct {
			// CompanyPageLogo - comment struct.
			CompanyPageLogo struct {
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage
			} `yaml:"company_page_logo"`
		} `yaml:"provider_account"`
		// FileStation - comment struct.
		FileStation struct {
			// ImageProxy - comment struct.
			ImageProxy struct {
				Host         string `yaml:"host" env:"APPX_IMAGE_HOST"`
				BasePath     string `yaml:"base_path"`
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage
			} `yaml:"image_proxy"`
		} `yaml:"file_station"`
	}

	// Validation - comment struct.
	Validation struct {
		Files struct {
			Json FileType `yaml:"json"`
		} `yaml:"files"`
		Images struct {
			Logo ImageType `yaml:"logo"`
		} `yaml:"images"`
		MimeTypes []mrlib.MimeType `yaml:"mime_types"`
	}

	// FileType - comment struct.
	FileType struct {
		MinSize                 uint64   `yaml:"min_size"`
		MaxSize                 uint64   `yaml:"max_size"`
		MaxFiles                uint32   `yaml:"max_files"`
		CheckRequestContentType bool     `yaml:"check_request_content_type"`
		Extensions              []string `yaml:"extensions"`
	}

	// ImageType - comment struct.
	ImageType struct {
		MaxWidth  uint64   `yaml:"max_width"`
		MaxHeight uint64   `yaml:"max_height"`
		CheckBody bool     `yaml:"check_body"`
		File      FileType `yaml:"file"`
	}

	// TaskSchedule - comment struct.
	TaskSchedule struct {
		ReloadSettings SchedulerTask `yaml:"reload_settings"`
		Mailer         struct {
			SendProcessor       MessageProcessor `yaml:"send_processor"`
			ChangeFromToRetry   SchedulerTask    `yaml:"change_from_to_retry"`
			CleanQueue          SchedulerTask    `yaml:"clean_queue"`
			SendRetryAttempts   uint32           `yaml:"send_retry_attempts"`
			SendDelayCorrection time.Duration    `yaml:"send_delay_correction"`
			ChangeQueueLimit    uint32           `yaml:"change_queue_limit"`
			ChangeRetryTimeout  time.Duration    `yaml:"change_retry_timeout"`
			ChangeRetryDelayed  time.Duration    `yaml:"change_retry_delayed"`
			CleanQueueLimit     uint32           `yaml:"clean_queue_limit"`
		} `yaml:"mailer"`
		Notifier struct {
			SendProcessor      MessageProcessor `yaml:"send_processor"`
			ChangeFromToRetry  SchedulerTask    `yaml:"change_from_to_retry"`
			CleanQueue         SchedulerTask    `yaml:"clean_queue"`
			SendRetryAttempts  uint32           `yaml:"send_retry_attempts"`
			ChangeQueueLimit   uint32           `yaml:"change_queue_limit"`
			ChangeRetryTimeout time.Duration    `yaml:"change_retry_timeout"`
			ChangeRetryDelayed time.Duration    `yaml:"change_retry_delayed"`
			CleanQueueLimit    uint32           `yaml:"clean_queue_limit"`
		} `yaml:"notifier"`
	}

	// MessageProcessor - comment struct.
	MessageProcessor struct {
		Caption           string        `yaml:"caption"`
		ReadyTimeout      time.Duration `yaml:"ready_timeout"`
		StartReadDelay    time.Duration `yaml:"start_read_delay"`
		ReadPeriod        time.Duration `yaml:"read_period"`
		CancelReadTimeout time.Duration `yaml:"cancel_read_timeout"`
		HandlerTimeout    time.Duration `yaml:"handler_timeout"`
		QueueSize         uint32        `yaml:"queue_size"`
		WorkersCount      uint16        `yaml:"workers_count"`
	}

	// SchedulerTask - comment struct.
	SchedulerTask struct {
		Caption string        `yaml:"caption"`
		Startup bool          `yaml:"startup"`
		Period  time.Duration `yaml:"period"`
		Timeout time.Duration `yaml:"timeout"`
	}
)
