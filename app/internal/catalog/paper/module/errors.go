package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

var (
	// ErrPaperNotFound - paper with ID not found.
	ErrPaperNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogPaperNotFound", mrerr.ErrorKindUser, "paper with ID={{ .id }} not found")

	// ErrPaperArticleAlreadyExists - paper article already exist.
	ErrPaperArticleAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogPaperArticleAlreadyExists", mrerr.ErrorKindUser, "paper article '{{ .name }}' already exist")
)

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrPaperNotFound,
		ErrPaperArticleAlreadyExists,
	}
}
