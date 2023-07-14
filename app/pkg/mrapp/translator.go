package mrapp

import "calc-user-data-back-adm/pkg/mrlang"

type Translator interface {
    GetLocale(langs ...mrlang.LangCode) Locale
}
