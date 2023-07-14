package mrlang

import "strings"

const LanguageByDefault = "en"

type (
	LangCode string // ISO 639 and regions

    ErrorMessage struct {
        Reason string `yaml:"reason"`
        Details []string `yaml:"details"`
    }
)

func CastToLangCodes(langs ...string) []LangCode {
    var langCodes []LangCode

    for _, lang := range langs {
        langCodes = append(langCodes, LangCode(lang))
    }

    return langCodes
}

func (em *ErrorMessage) DetailsToString() string {
    switch len(em.Details) {
    case 0:
        return ""
    case 1:
        return em.Details[0]
    }

    return "- " + strings.Join(em.Details, "\n- ")
}
