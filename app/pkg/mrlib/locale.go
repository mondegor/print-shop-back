package mrlib

import (
    "fmt"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrlang"

    "github.com/ilyakaznacheev/cleanenv"
)

type Locale struct {
    code mrlang.LangCode
    Messages map[string]string `yaml:"messages"`
    Errors map[string]mrlang.ErrorMessage `yaml:"errors"`
}

// Make sure the Locale conforms with the mrapp.Locale interface
var _ mrapp.Locale = (*Locale)(nil)

func NewLocale(code mrlang.LangCode, filePath string) (*Locale, error) {
    locale := &Locale{
        code: code,
    }

    if err := cleanenv.ReadConfig(filePath, locale); err != nil {
        return nil, fmt.Errorf("while reading locale '%s', error '%s' occurred", filePath, err)
    }

    return locale, nil
}

func (l Locale) GetLang() mrlang.LangCode {
    return l.code
}

func (l Locale) GetMessage(code string, defaultMessage string, args ...mrlang.NamedArg) string {
    value, ok := l.Messages[code]

    if !ok {
        value = defaultMessage
    }

    if len(args) > 0 {
        value = mrlang.RenderMessage(value, args)
    }

    return value
}

func (l Locale) GetError(code string, defaultMessage string, args ...mrlang.NamedArg) mrlang.ErrorMessage {
    value, ok := l.Errors[code]

    if !ok {
        value = mrlang.ErrorMessage{Reason: defaultMessage}
    }

    if len(args) > 0 {
        value.Reason = mrlang.RenderMessage(value.Reason, args)

        for i := 0; i < len(value.Details); i++ {
            value.Details[i] = mrlang.RenderMessage(value.Details[i], args)
        }
    }

    return value
}
