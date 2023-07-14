package mrerr

import (
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrlang"
)

func (e *AppError) GetUserMessage(locale mrapp.Locale) mrlang.ErrorMessage {
    if e.kind != ErrorKindInternal {
        return locale.GetError(string(e.code), e.message, e.getNamedArgs()...)
    }

    return locale.GetError(ErrorCodeInternal, ErrorCodeInternal)
}

func (e *AppError) renderMessage() []byte {
    if len(e.argsNames) == 0 {
        return []byte(e.message)
    }

    return []byte(mrlang.RenderMessage(e.message, e.getNamedArgs()))
}

func (e *AppError) getNamedArgs() []mrlang.NamedArg {
    var namedArgs []mrlang.NamedArg

    for i, argName := range e.argsNames {
        namedArgs = append(namedArgs, mrlang.NewArg(argName, e.args[i]))
    }

    return namedArgs
}

func getMessageArgsNames(message string) []string {
    return mrlang.GetMessageArgsNames(message)
}
