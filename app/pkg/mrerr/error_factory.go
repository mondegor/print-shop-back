package mrerr

import (
    "errors"
    "fmt"
    "runtime"

    "github.com/google/uuid"
)

type (
    AppErrorFactory struct {
        code ErrorCode
        kind ErrorKind
        message string
        argsNames []string
        callerSkip int
    }
)

func NewFactory(code ErrorCode, kind ErrorKind, message string) *AppErrorFactory {
    return &AppErrorFactory{
        code: code,
        kind: kind,
        message: message,
        argsNames: getMessageArgsNames(message),
    }
}

func (e *AppErrorFactory) New(args ...any) *AppError {
    return e.new(nil, args)
}

func (e *AppErrorFactory) Wrap(err error, args ...any) *AppError {
    if err == nil {
        panic("error is nil, wrapping is not necessary")
    }

    return e.new(err, args)
}

func (e *AppErrorFactory) Caller(skip int) *AppErrorFactory {
    return &AppErrorFactory{
        code: e.code,
        kind: e.kind,
        message: e.message,
        argsNames: e.argsNames,
        callerSkip: e.callerSkip + skip,
    }
}

// Is - see: AppError::Is
func (e *AppErrorFactory) Is(err error) bool {
    return errors.Is(err, &AppError{code: e.code})
}

func (e *AppErrorFactory) new(err error, args []any) *AppError {
    newErr := &AppError{
        code: e.code,
        kind: e.kind,
        message: e.message,
        argsNames: e.argsNames,
        args: args,
        err: err,
    }

    e.init(newErr)

    return newErr
}

func (e *AppErrorFactory) init(newErr *AppError) {
    newErr.setErrorIfArgsNotEqual(4)

    if newErr.err != nil {
        appError, ok := newErr.err.(*AppError)

        // raising to the top
        if ok && appError.eventId != nil {
            newErr.eventId = appError.eventId
            appError.eventId = nil
            return
        }
    }

    if e.kind != ErrorKindInternal && e.kind != ErrorKindSystem {
        return
    }

    _, file, line, ok := runtime.Caller(e.callerSkip + 3)

    if ok {
        if file == "" {
            file = "???"
        }

        newErr.file = new(string)
        *newErr.file = file
        newErr.line = line
    }

    //if e.eventId == nil {
    //    e.eventId = (*string)(sentry.CaptureException(e))
    //}

    if newErr.eventId == nil {
        newErr.eventId = new(string)
        *newErr.eventId = fmt.Sprintf("ER-%s", uuid.New().String())
    }
}
