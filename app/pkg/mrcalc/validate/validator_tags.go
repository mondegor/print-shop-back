package validate

import (
	"regexp"
)

var (
	regexp2dSize = regexp.MustCompile(`^[0-9]+x[0-9]+$`)
	regexp3dSize = regexp.MustCompile(`^[0-9]+x[0-9]+x[0-9]+$`)
)

// Size2d - comment func.
func Size2d(value string) bool {
	return regexp2dSize.MatchString(value)
}

// Size3d - comment func.
func Size3d(value string) bool {
	return regexp3dSize.MatchString(value)
}
