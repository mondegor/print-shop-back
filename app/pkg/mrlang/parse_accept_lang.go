package mrlang

import (
    "regexp"
    "strings"
)

const maxAcceptLanguageLen = 256

var regexpAcceptLanguage = regexp.MustCompile("^[a-z]{2}(-[a-zA-Z0-9-]+)?$")

// ParseAcceptLanguage
// Sample Accept-Language: ru;q=0.9, fr-CH, fr;q=0.8, en;q=0.7, *;q=0.5
func ParseAcceptLanguage(s string) []LangCode {
    length := len(s)

    if length > 0 && length <= maxAcceptLanguageLen {
        var langs []LangCode
        var keys map[string]bool

        for _, lang := range strings.Split(strings.ToLower(s), ",") {
            if index := strings.Index(lang, ";"); index >= 0 {
                lang = lang[:index]
            }

            lang = strings.TrimSpace(lang)

            if !regexpAcceptLanguage.MatchString(lang) {
                continue
            }

            if keys == nil {
                keys = make(map[string]bool, 2)
            } else {
                if _, ok := keys[lang]; ok {
                    continue
                }
            }

            langs = append(langs, LangCode(lang))
            keys[lang] = true
        }

        if len(langs) > 0 {
            return langs
        }
    }

    return []LangCode{LanguageByDefault}
}
