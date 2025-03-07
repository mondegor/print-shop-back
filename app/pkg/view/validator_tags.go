package view

import (
	"regexp"
)

var (
	regexp2dSize = regexp.MustCompile(`^[0-9]+x[0-9]+$`)
	regexp3dSize = regexp.MustCompile(`^[0-9]+x[0-9]+x[0-9]+$`)
)

// Validate2dSize - comment func.
func Validate2dSize(value string) bool {
	return regexp2dSize.MatchString(value)
}

// Validate3dSize - comment func.
func Validate3dSize(value string) bool {
	return regexp3dSize.MatchString(value)
}
