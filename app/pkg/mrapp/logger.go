package mrapp

import (
    "log"
    "strings"
)

type LogLevel uint32

const (
    LogFatalLevel LogLevel = iota
    LogErrorLevel
    LogWarnLevel
    LogInfoLevel
    LogDebugLevel
)

type Logger interface {
    WithContext(name string) Logger
    Fatal(message any, args ...any)
    Error(message any, args ...any)
    Warn(message string, args ...any)
    Info(message string, args ...any)
    Debug(message any, args ...any)
    Event(message string, args ...any)
}

func ParseLogLevel(level string) LogLevel {
    switch strings.ToLower(level) {
    case "fatal":
        return LogFatalLevel

    case "error":
        return LogErrorLevel

    case "warn", "warning":
        return LogWarnLevel

    case "info":
        return LogInfoLevel

    case "debug":
        return LogDebugLevel
    }

    log.Fatalf("Log level '%s' is incorrect", level)

    return LogFatalLevel
}

