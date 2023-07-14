package mrapp

import "calc-user-data-back-adm/pkg/mrlang"

type LocaleStub struct { }

// Make sure the LocaleStub conforms with the Locale interface
var _ Locale = (*LocaleStub)(nil)

func NewLocaleStub() *LocaleStub {
    return &LocaleStub{}
}

func (l *LocaleStub) GetLang() mrlang.LangCode {
    return mrlang.LanguageByDefault
}

func (l *LocaleStub) GetMessage(code string, defaultMessage string, args ...mrlang.NamedArg) string {
    if len(args) > 0 {
        defaultMessage = mrlang.RenderMessage(defaultMessage, args)
    }

    return defaultMessage
}

func (l *LocaleStub) GetError(code string, defaultMessage string, args ...mrlang.NamedArg) mrlang.ErrorMessage {
    value := mrlang.ErrorMessage{Reason: defaultMessage}

    if len(args) > 0 {
        value.Reason = mrlang.RenderMessage(value.Reason, args)

        for i := 0; i < len(value.Details); i++ {
            value.Details[i] = mrlang.RenderMessage(value.Details[i], args)
        }
    }

    return value
}
