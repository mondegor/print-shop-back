package mrlib

import (
    "fmt"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrlang"
)

type (
    langMap map[mrlang.LangCode]*Locale

    Translator struct {
        logger mrapp.Logger
        langs langMap
        defaultLocale mrapp.Locale
    }

    TranslatorOptions struct {
        DirPath string
        FileType string
        LangCodes []mrlang.LangCode
    }
)

// Make sure the Translator conforms with the mrapp.Translator interface
var _ mrapp.Translator = (*Translator)(nil)

func NewTranslator(logger mrapp.Logger, opt TranslatorOptions) *Translator {
    langs := langMap{}
    var defaultLocale mrapp.Locale

    if len(opt.LangCodes) == 0 {
        defaultLocale = mrapp.NewLocaleStub()
        logger.Warn("opt.LangCodes is empty")
    } else {
        for i, code := range opt.LangCodes {
            locale, err := NewLocale(code, fmt.Sprintf("%s/%s.%s", opt.DirPath, code, opt.FileType))

            if err != nil {
                logger.Error(err)
                continue
            }

            langs[code] = locale

            if i == 0 {
                defaultLocale = locale
            }
        }
    }

    return &Translator{
        logger: logger,
        langs: langs,
        defaultLocale: defaultLocale,
    }
}

func (t Translator) GetLocale(langs ...mrlang.LangCode) mrapp.Locale {
    for _, lang := range langs {
        if locale, ok := t.langs[lang]; ok {
            return locale
        }
    }

    return t.defaultLocale
}
