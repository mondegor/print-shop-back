package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

var (
	// ErrLaminateNotFound - laminate with ID not found.
	ErrLaminateNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogLaminateNotFound", mrerr.ErrorKindUser, "laminate with ID={{ .id }} not found")

	// ErrLaminateArticleAlreadyExists - laminate article already exist.
	ErrLaminateArticleAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogLaminateArticleAlreadyExists", mrerr.ErrorKindUser, "laminate article '{{ .name }}' already exist")
)

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrLaminateNotFound,
		ErrLaminateArticleAlreadyExists,
	}
}
