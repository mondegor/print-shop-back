package mrapp

import "print-shop-back/pkg/mrlang"

type Translator interface {
    GetLocale(langs ...mrlang.LangCode) Locale
}
