package factory

import (
	"errors"
	"flag"
	"os"
)

const (
	dotEnvPathName = "APPX_ENV_PATH"
)

type (
	// AppArgs - разобранные аргументы, которые передаются из командной строки или тестовой среды.
	AppArgs struct {
		WorkDir     string // путь к рабочей директории приложения
		ConfigPath  string // путь к файлу конфигурации приложения
		DotEnvPath  string // путь к .env файлу (переменные из этого файла более приоритетны переменных из ConfigPath)
		Environment string // внешнее окружение: local, dev, test, prod
		LogLevel    string // уровень логирования: info, warn, error, fatal, debug, trace
	}
)

// ErrParseArgsHelp is the error returned if the -help or -h flag is invoked
// but no such flag is defined.
var ErrParseArgsHelp = flag.ErrHelp

// ParseAppArgs - возвращает разобранные аргументы или ошибку, если разбор не был успешным.
// Также может быть возвращено ErrParseArgsHelp в качестве ошибки.
func ParseAppArgs(args []string) (appArgs AppArgs, err error) {
	appArgs.DotEnvPath = os.Getenv(dotEnvPathName)

	if len(args) < 2 {
		return appArgs, nil
	}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)

	fs.StringVar(&appArgs.WorkDir, "work-dir", "", "Path to the application work dir")
	fs.StringVar(&appArgs.ConfigPath, "config-path", "", "Path to the application config file")
	fs.StringVar(&appArgs.Environment, "environment", "", "App environment (local, dev, test, prod)")
	fs.StringVar(&appArgs.LogLevel, "log-level", "", "Logging level (info, warn, error, fatal, debug, trace)")

	if err = fs.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return AppArgs{}, ErrParseArgsHelp
		}

		return AppArgs{}, err
	}

	return appArgs, nil
}
