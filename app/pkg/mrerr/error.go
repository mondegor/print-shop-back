package mrerr

import (
    "errors"
    "fmt"
)

const (
    ErrorKindInternal ErrorKind = iota
    ErrorKindInternalNotice
    ErrorKindSystem
    ErrorKindUser

    ErrorCodeInternal = "errInternal"
)

type (
    ErrorCode string
    ErrorKind int32

    AppError struct {
        code ErrorCode
        kind ErrorKind
        message string
        argsNames []string
        args []any
        err error
        file *string
        line int
        eventId *string
    }
)

func New(code ErrorCode, message string, args ...any) *AppError {
    newErr := &AppError{
        code: code,
        kind: ErrorKindUser,
        message: message,
        argsNames: getMessageArgsNames(message),
        args: args,
    }

    newErr.setErrorIfArgsNotEqual()

    return newErr
}

func (e *AppError) setErrorIfArgsNotEqual() {
    if len(e.argsNames) == len(e.args) {
        return
    }

    var argsErrorFactory *AppErrorFactory

    if len(e.argsNames) > len(e.args) {
        argsErrorFactory = ErrInternalMessageNotEnoughArguments
    } else {
        argsErrorFactory = ErrInternalMessageTooManyArguments
    }

    if e.err == nil {
        e.err = argsErrorFactory.New(e.message)
        return
    }

    // infinite loop protection
    if !errors.Is(e.err, &AppError{code: argsErrorFactory.code}) {
        e.err = argsErrorFactory.Wrap(e.err, e.message)
    }
}

func (e *AppError) Is(err error) bool {
    if v, ok := err.(*AppError); ok && e.code == v.code {
        return true
    }

    return false
}

func (e *AppError) Error() string {
    var buf []byte

    if e.eventId != nil {
        buf = append(buf, '[')
        buf = append(buf, *e.eventId...)
        buf = append(buf, ']', ' ')
    }

    buf = append(buf, e.renderMessage()...)

    if e.file != nil {
        buf = append(buf, fmt.Sprintf(" in %s:%d", *e.file, e.line)...)
    }

    if e.err != nil {
        buf = append(buf, ';', ' ')
        buf = append(buf, e.err.Error()...)
    }

    return string(buf)
}

func (e *AppError) Unwrap() error {
    return e.err
}

func (e *AppError) EventId() string {
    if e.eventId == nil {
        return ""
    }

    return *e.eventId
}
