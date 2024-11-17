package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrBoxNotFound - box with ID not found.
	ErrBoxNotFound = mrerr.NewProto(
		"catalog.errBoxNotFound", mrerr.ErrorKindUser, "box with ID={{ .id }} not found")

	// ErrBoxArticleAlreadyExists - box article already exist.
	ErrBoxArticleAlreadyExists = mrerr.NewProto(
		"catalog.errBoxArticleAlreadyExists", mrerr.ErrorKindUser, "box article '{{ .name }}' already exist")
)
