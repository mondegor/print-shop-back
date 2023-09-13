package view

import (
    "regexp"
)

var (
    regexpArticle = regexp.MustCompile(`^\S+$`)
    regexpVariable = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
)

func ValidateArticle(value string) bool {
    return regexpArticle.MatchString(value)
}

func ValidateVariable(value string) bool {
    return regexpVariable.MatchString(value)
}
