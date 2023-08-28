package config

import (
    "log"
    "os"

    "github.com/ilyakaznacheev/cleanenv"
)

type (
    Config struct {
        AppPath string `yaml:"app_path"`
        Debug bool `yaml:"debug" env:"APPX_DEBUG"`
        Listen `yaml:"listen"`
        Log `yaml:"logger"`
        Storage `yaml:"storage"`
        // Sentry `yaml:"sentry"`
        Cors `yaml:"cors"`
        Translation `yaml:"translation"`
    }

    Listen struct {
        Type string `yaml:"type" env:"APPX_SERVICE_LISTEN_TYPE"`
        SockName string `yaml:"sock_name" env:"APPX_SERVICE_LISTEN_SOCK"`
        BindIP string `yaml:"bind_ip" env:"APPX_SERVICE_BIND"`
        Port string `yaml:"port" env:"APPX_SERVICE_PORT"`
    }

    Log struct {
        Level string `yaml:"level" env:"APPX_LOG_LEVEL"`
        NoColor bool `yaml:"no_color" env:"APPX_NO_COLOR"`
    }

    Storage struct {
        Host string `yaml:"host" env:"APPX_DB_HOST"`
        Port string `yaml:"port" env:"APPX_DB_PORT"`
        Username string `yaml:"username" env:"APPX_DB_USER"`
        Password string `yaml:"password" env:"APPX_DB_PASSWORD"`
        Database string `yaml:"database" env:"APPX_DB_NAME"`
        Timeout int `yaml:"timeout"` // in sec
    }

    //Redis struct {
    //    Host string `yaml:"host" env:"APPX_REDIS_HOST"`
    //    Port string `yaml:"port" env:"APPX_REDIS_PORT"`
    //    Password string `yaml:"password" env:"APPX_REDIS_PASSWORD"`
    //    Timeout int `yaml:"timeout"` // in sec
    //}

    //Sentry struct {
    //    Use bool `yaml:"use" env:"APPX_SENTRY_USE"`
    //    Dsn string `yaml:"dsn" env:"APPX_SENTRY_DSN"`
    //}

    Cors struct {
        AllowedOrigins []string `yaml:"allowed_origins"`
        AllowedMethods []string `yaml:"allowed_methods"`
        AllowedHeaders []string `yaml:"allowed_headers"`
        ExposedHeaders []string `yaml:"exposed_headers"`
        AllowCredentials bool `yaml:"allow_credentials"`
    }

    Translation struct {
        DirPath string `yaml:"dir_path"`
        FileType string `yaml:"file_type"`
        LangCodes []string `yaml:"lang_codes"`
    }
)

func New(filePath string) *Config {
    cfg := &Config{}
    err := cleanenv.ReadConfig(filePath, cfg)

    if err != nil {
        log.Fatalf("While reading config '%s', error '%s' occurred", filePath, err)
    }

    err = cleanenv.ReadEnv(cfg)

    if err != nil {
        log.Fatal(err)
    }

    if cfg.AppPath == "" {
        cfg.AppPath = os.Args[0]
    }

    return cfg
}
