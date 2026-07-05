package log

import (
	"github.com/mondegor/go-core/mrlog"
)

type (
	// Logger - интерфейс логирования сообщений и ошибок с использованием контекста.
	Logger = mrlog.Logger
)

// WithAttrs - возвращает логгер с прикреплёнными к нему атрибутами.
func WithAttrs(l Logger, attrs ...any) Logger {
	return mrlog.WithAttrs(l, attrs...)
}

// // DebugEnabled - сообщает логирует ли указанный логгер сообщения уровня level.Debug.
// // Если logger = nil, то будет возвращено false.
// func DebugEnabled(l Logger) bool {
// 	return mrlog.Enabled(l, level.Debug)
// }

// // InfoEnabled - сообщает логирует ли указанный логгер сообщения уровня level.Info.
// // Если logger = nil, то будет возвращено false.
// func InfoEnabled(l Logger) bool {
// 	return mrlog.Enabled(l, level.Info)
// }

// Debug - логирует сообщения на уровне level.Debug без использования контекста.
func Debug(l Logger, msg string, args ...any) {
	mrlog.Debug(l, msg, args...)
}

// DebugFunc - логирует сообщения на уровне level.Debug без использования контекста с их отложенным созданием сообщения.
func DebugFunc(l Logger, createMsg func() string, args ...any) {
	mrlog.DebugFunc(l, createMsg, args...)
}

// Info - логирует сообщения на уровне level.Info без использования контекста.
func Info(l Logger, msg string, args ...any) {
	mrlog.Info(l, msg, args...)
}

// Warn - логирует сообщения на уровне level.Warn без использования контекста.
func Warn(l Logger, msg string, args ...any) {
	mrlog.Warn(l, msg, args...)
}

// // Error - логирует сообщения на уровне level.Error без использования контекста.
// func Error(l Logger, msg string, args ...any) {
// 	mrlog.Error(l, msg, args...)
// }

// FatalError - логирует сообщения на уровне level.Fatal без использования контекста и останавливает приложение.
func FatalError(l Logger, msg string, args ...any) {
	mrlog.FatalError(l, msg, args...)
}

// Fatal is equivalent to [Print] followed by a call to [os.Exit](1).
func Fatal(v ...any) {
	mrlog.Fatal(v...)
}

// NopLogger - создаёт объект Logger, который ничего не делает.
func NopLogger() Logger {
	return mrlog.NopLogger()
}
