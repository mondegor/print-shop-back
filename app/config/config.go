package config

import (
	"io"
	"time"

	authcfg "github.com/mondegor/go-components/wire/mrauth/config"
	modelcfg "github.com/mondegor/go-sysmess/mrmodel/config"
	"github.com/mondegor/go-sysmess/util/mime"
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
		App                   `yaml:"app"`
		Debugging             `yaml:"debugging"`
		Log                   `yaml:"logger"`
		Trace                 `yaml:"tracer"`
		Sentry                `yaml:"sentry"`
		Servers               `yaml:"servers"`
		Storage               `yaml:"storage"`
		Redis                 `yaml:"redis"`
		FileSystem            `yaml:"file_system"` // disabled
		S3                    `yaml:"s3"`
		FileProviders         `yaml:"file_providers"`
		Cors                  `yaml:"cors"`
		Localization          `yaml:"localization"`
		Senders               `yaml:"senders"`
		authcfg.AccessControl `yaml:"access_control"`
		ModulesSettings       `yaml:"modules_settings"`
		Validation            `yaml:"validation"`
		TaskSchedule          `yaml:"task_schedule"`
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
		Debug                  bool `yaml:"debug" env:"APPX_DEBUG"`
		authcfg.AuthorizedUser `yaml:"authorized_user"`
		UnexpectedHttpStatus   uint16 `yaml:"unexpected_http_status"`
		ErrorCaller            `yaml:"error_caller"`
	}

	// ErrorCaller - comment struct.
	ErrorCaller struct {
		IsEnabled   bool     `yaml:"is_enabled" env:"APPX_ERR_CALLER_IS_ENABLED"`
		Depth       uint8    `yaml:"depth" env:"APPX_ERR_CALLER_DEPTH"`
		ShowFunc    bool     `yaml:"show_func" env:"APPX_ERR_CALLER_SHOW_FUNC"`
		UpperBounds []string `yaml:"upper_bounds"`
	}

	// Log - comment struct.
	Log struct {
		Level      string `yaml:"level" env:"APPX_LOG_LEVEL"`
		TimeFormat string `yaml:"time_format" env:"APPX_LOG_TIME_FORMAT"`
		JsonFormat bool   `yaml:"json_format" env:"APPX_LOG_JSON"`
		ColorMode  bool   `yaml:"color_mode" env:"APPX_LOG_COLOR"`
	}

	Trace struct {
		IsEnabled bool `yaml:"is_enabled" env:"APPX_TRACE_IS_ENABLED"`
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
		MaxPoolSize     uint16        `yaml:"max_pool_size" env:"APPX_DB_MAX_POOL_SIZE"`
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

	// FileSystem - comment struct (disabled).
	FileSystem struct {
		DirMode    uint32 `yaml:"dir_mode" env:"APPX_FILESYSTEM_DIR_MODE"`
		CreateDirs bool   `yaml:"create_dirs" env:"APPX_FILESYSTEM_CREATE_DIRS"`
	}

	// S3 - comment struct.
	S3 struct {
		Host          string `yaml:"host" env:"APPX_S3_HOST"`
		Port          string `yaml:"port" env:"APPX_S3_PORT"`
		UseSSL        bool   `yaml:"use_ssl" env:"APPX_S3_USESSL"`
		Username      string `yaml:"username" env:"APPX_S3_USER"`
		Password      string `yaml:"password" env:"APPX_S3_PASSWORD"`
		CreateBuckets bool   `yaml:"create_buckets" env:"APPX_S3_CREATE_BUCKETS"`
	}

	// FileProviders - comment struct.
	FileProviders struct {
		// ImageStorage - comment struct.
		ImageStorage struct {
			Name       string `yaml:"name"`
			BucketName string `yaml:"bucket_name" env:"APPX_IMAGESTORAGE_BUCKET"` // S3
			RootDir    string `yaml:"root_dir" env:"APPX_IMAGESTORAGE_ROOT_DIR"`  // FileSystem (disabled)
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

	// Localization - comment struct.
	Localization struct {
		LangURLParam string   `yaml:"lang_url_param"`
		Languages    []string `yaml:"languages" env:"APPX_LOCALIZATION_LANGS"` // items by "," separated
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

	// ModulesSettings - comment struct.
	ModulesSettings struct {
		// General - comment struct.
		General struct {
			PageSizeMax     uint16 `yaml:"page_size_max"`
			PageSizeDefault uint16 `yaml:"page_size_default"`
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
			Json modelcfg.FileType `yaml:"json"`
		} `yaml:"files"`
		Images struct {
			Logo modelcfg.ImageType `yaml:"logo"`
		} `yaml:"images"`
		MimeTypes []mime.Type `yaml:"mime_types"`
	}

	// TaskSchedule - comment struct.
	TaskSchedule struct {
		Settings struct {
			// Caption        string        `yaml:"caption"`
			ReloadSettings     SchedulerTask `yaml:"reload_settings"`
			DefaultPeriodRatio float64       `yaml:"default_default_period_ratio"`
		} `yaml:"settings"`
		Auth struct {
			// Caption           string        `yaml:"caption"`
			CleanRecords      SchedulerTask `yaml:"clean_records"`
			CleanRecordsLimit uint32        `yaml:"clean_records_limit"`
			LogsLifeTime      time.Duration `yaml:"logs_life_time"`
			UserStat          struct {
				RequestCollector MessageCollector `yaml:"request_collector"`
			} `yaml:"user_stat"`
		} `yaml:"auth"`
		Mailer struct {
			// Caption             string           `yaml:"caption"`
			MessageProcessor     MessageProcessor `yaml:"message_processor"`
			ChangeFromToRetry    SchedulerTask    `yaml:"change_from_to_retry"`
			CleanQueue           SchedulerTask    `yaml:"clean_queue"`
			SendRetryAttempts    uint8            `yaml:"send_retry_attempts"`
			SendDelayCorrection  time.Duration    `yaml:"send_delay_correction"`
			ChangeQueueBatchSize uint32           `yaml:"change_queue_batch_size"`
			ChangeRetryTimeout   time.Duration    `yaml:"change_retry_timeout"`
			ChangeRetryDelayed   time.Duration    `yaml:"change_retry_delayed"`
			CleanQueueBatchSize  uint32           `yaml:"clean_queue_batch_size"`
		} `yaml:"mailer"`
		Notifier struct {
			// Caption            string           `yaml:"caption"`
			NoticeProcessor      MessageProcessor `yaml:"notice_processor"`
			ChangeFromToRetry    SchedulerTask    `yaml:"change_from_to_retry"`
			CleanQueue           SchedulerTask    `yaml:"clean_queue"`
			SendRetryAttempts    uint8            `yaml:"send_retry_attempts"`
			ChangeQueueBatchSize uint32           `yaml:"change_queue_batch_size"`
			ChangeRetryTimeout   time.Duration    `yaml:"change_retry_timeout"`
			ChangeRetryDelayed   time.Duration    `yaml:"change_retry_delayed"`
			CleanQueueBatchSize  uint32           `yaml:"clean_queue_batch_size"`
		} `yaml:"notifier"`
	}

	// MessageCollector - comment struct.
	MessageCollector struct {
		// Caption              string        `yaml:"caption"`
		ReadyTimeout   time.Duration `yaml:"ready_timeout"`
		FlushPeriod    time.Duration `yaml:"flush_period"`
		HandlerTimeout time.Duration `yaml:"handler_timeout"`
		BatchSize      uint32        `yaml:"batch_size"`
		WorkersCount   uint8         `yaml:"workers_count"`
	}

	// MessageProcessor - comment struct.
	MessageProcessor struct {
		// Caption              string        `yaml:"caption"`
		ReadyTimeout         time.Duration `yaml:"ready_timeout"`
		ReadPeriod           time.Duration `yaml:"read_period"`
		ConsumerReadTimeout  time.Duration `yaml:"consumer_read_timeout"`
		ConsumerWriteTimeout time.Duration `yaml:"consumer_write_timeout"`
		HandlerTimeout       time.Duration `yaml:"handler_timeout"`
		QueueSize            uint16        `yaml:"queue_size"`
		WorkersCount         uint8         `yaml:"workers_count"`
		NotificationChannel  string        `yaml:"notification_channel,omitempty"`
	}

	// SchedulerTask - comment struct.
	SchedulerTask struct {
		// Caption             string        `yaml:"caption"`
		// Startup             bool          `yaml:"startup"`
		Period              time.Duration `yaml:"period"`
		Timeout             time.Duration `yaml:"timeout"`
		NotificationChannel string        `yaml:"notification_channel,omitempty"`
	}
)
