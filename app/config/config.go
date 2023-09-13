package config

import (
    "fmt"
    "os"

    "github.com/ilyakaznacheev/cleanenv"
)

type (
    Config struct {
        AppPath string `yaml:"app_path"`
        Debug bool `yaml:"debug" env:"APPX_DEBUG"`
        Server `yaml:"server"`
        Listen `yaml:"listen"`
        Log `yaml:"logger"`
        Storage `yaml:"storage"`
        Cors `yaml:"cors"`
        Translation `yaml:"translation"`
    }

    Server struct {
        ReadTimeout int32 `yaml:"readTimeout"`
        WriteTimeout int32 `yaml:"writeTimeout"`
        ShutdownTimeout int32 `yaml:"shutdownTimeout"`
    }

    Listen struct {
        Type string `yaml:"type" env:"APPX_SERVICE_LISTEN_TYPE"`
        SockName string `yaml:"sock_name" env:"APPX_SERVICE_LISTEN_SOCK"`
        BindIP string `yaml:"bind_ip" env:"APPX_SERVICE_BIND"`
        Port string `yaml:"port" env:"APPX_SERVICE_PORT"`
    }

    Log struct {
        Level string `yaml:"level" env:"APPX_LOG_LEVEL"`
    }

    Storage struct {
        Host string `yaml:"host" env:"APPX_DB_HOST"`
        Port string `yaml:"port" env:"APPX_DB_PORT"`
        Username string `yaml:"username" env:"APPX_DB_USER"`
        Password string `yaml:"password" env:"APPX_DB_PASSWORD"`
        Database string `yaml:"database" env:"APPX_DB_NAME"`
        MaxPoolSize int32 `yaml:"maxPoolSize" env:"APPX_DB_MAX_POOL_SIZE"`
        Timeout int32 `yaml:"timeout"` // in sec
    }

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

func New(filePath string) (*Config, error) {
    cfg := &Config{}
    err := cleanenv.ReadConfig(filePath, cfg)

    if err != nil {
        return nil, fmt.Errorf("while reading config '%s', error '%s' occurred", filePath, err)
    }

    err = cleanenv.ReadEnv(cfg)

    if err != nil {
        return nil, err
    }

    if cfg.AppPath == "" {
        cfg.AppPath = os.Args[0]
    }

    return cfg, nil
}
