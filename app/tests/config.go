package tests

import (
	"path"
	"runtime"
	"sync"
)

// AppWorkDir - возвращает рабочую директорию приложения.
func AppWorkDir() string {
	var (
		once    sync.Once
		workDir string
	)

	once.Do(func() {
		_, p, _, _ := runtime.Caller(0)
		workDir = path.Join(path.Dir(p), "..") // up from '.../tests'
	})

	return workDir
}

// AppMigrationsDir - возвращает директорию с миграциями приложения.
func AppMigrationsDir() string {
	return AppWorkDir() + "/migrations"
}

// AppDotEnvPathForTests - возвращает директорию с миграциями приложения.
func AppDotEnvPathForTests() string {
	return AppWorkDir() + "/tests/.env"
}

// DBSchemas - возвращает массив схем БД, с которыми работает приложение.
func DBSchemas() []string {
	return []string{
		"printshop_calculation",
		"printshop_catalog",
		"printshop_controls",
		"printshop_dictionaries",
		"printshop_global",
		"printshop_providers",
	}
}

// ExcludedDBTables - возвращает массив таблиц БД, которые не должны меняться при тестировании.
func ExcludedDBTables() []string {
	return nil
}
