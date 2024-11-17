package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrPaperNotFound - paper with ID not found.
	ErrPaperNotFound = mrerr.NewProto(
		"catalog.errPaperNotFound", mrerr.ErrorKindUser, "paper with ID={{ .id }} not found")

	// ErrPaperArticleAlreadyExists - paper article already exist.
	ErrPaperArticleAlreadyExists = mrerr.NewProto(
		"catalog.errPaperArticleAlreadyExists", mrerr.ErrorKindUser, "paper article '{{ .name }}' already exist")
)
