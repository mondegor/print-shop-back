package config

import (
	"errors"
	"flag"
)

type (
	// CmdArgs - разобранные аргументы, которые передаются из командной строки или тестовой среды.
	// Эти аргументы более приоритетны аналогичным, определенным в конфигурации или переданных через .env файл.
	CmdArgs struct {
		WorkDir     string // путь к рабочей директории приложения
		Environment string // внешнее окружение: local, dev, test, prod
		DotEnvPath  string // путь к .env файлу (переменные из этого файла более приоритетны переменных из ConfigPath)
		LogLevel    string // уровень логирования: info, warn, error, fatal, debug, trace
	}
)

// ErrParseArgsHelp is the error returned if the -help or -h flag is invoked
// but no such flag is defined.
var ErrParseArgsHelp = flag.ErrHelp

// ParseCmdArgs - возвращает разобранные аргументы или ошибку, если разбор не был успешным.
// Также может быть возвращено ErrParseArgsHelp в качестве ошибки.
func ParseCmdArgs(args []string) (cmdArgs CmdArgs, err error) {
	if len(args) < 2 {
		return cmdArgs, nil
	}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)

	fs.StringVar(&cmdArgs.WorkDir, "work-dir", "", "Path to the application work dir")
	fs.StringVar(&cmdArgs.Environment, "environment", "", "App environment (local, dev, test, prod)")
	fs.StringVar(&cmdArgs.DotEnvPath, "dot-env-path", "", "Path to the .env file")
	fs.StringVar(&cmdArgs.LogLevel, "log-level", "", "Logging level (info, warn, error, fatal, debug, trace)")

	if err = fs.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return CmdArgs{}, ErrParseArgsHelp
		}

		return CmdArgs{}, err
	}

	return cmdArgs, nil
}
