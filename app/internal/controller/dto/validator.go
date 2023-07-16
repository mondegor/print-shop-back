package dto

import (
    "regexp"

    "github.com/go-playground/validator/v10"
)

var (
    regexpArticle = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
    regexpVariable = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
)

func ValidateArticle() any {
    return func (fl validator.FieldLevel) bool {
        return regexpArticle.MatchString(fl.Field().String())
    }
}

func ValidateVariable() any {
    return func (fl validator.FieldLevel) bool {
        return regexpVariable.MatchString(fl.Field().String())
    }
}
