package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

var (
	// ErrBoxNotFound - box with ID not found.
	ErrBoxNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogBoxNotFound", mrerr.ErrorKindUser, "box with ID={{ .id }} not found")

	// ErrBoxArticleAlreadyExists - box article already exist.
	ErrBoxArticleAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogBoxArticleAlreadyExists", mrerr.ErrorKindUser, "box article '{{ .name }}' already exist")
)

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrBoxNotFound,
		ErrBoxArticleAlreadyExists,
	}
}
