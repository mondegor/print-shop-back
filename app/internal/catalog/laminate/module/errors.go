package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrLaminateNotFound - laminate with ID not found.
	ErrLaminateNotFound = mrerr.NewProto(
		"catalog.errLaminateNotFound", mrerr.ErrorKindUser, "laminate with ID={{ .id }} not found")

	// ErrLaminateArticleAlreadyExists - laminate article already exist.
	ErrLaminateArticleAlreadyExists = mrerr.NewProto(
		"catalog.errLaminateArticleAlreadyExists", mrerr.ErrorKindUser, "laminate article '{{ .name }}' already exist")
)
