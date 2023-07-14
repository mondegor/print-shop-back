package mrapp

import "calc-user-data-back-adm/pkg/mrlang"

type Locale interface {
    GetLang() mrlang.LangCode
    GetMessage(code string, defaultMessage string, args ...mrlang.NamedArg) string
    GetError(code string, defaultMessage string, args ...mrlang.NamedArg) mrlang.ErrorMessage
}
