package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

// ErrCompanyPageRewriteNameAlreadyExists - rewrite name already exists.
var ErrCompanyPageRewriteNameAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
	"errProviderAccountsCompanyPageRewriteNameAlreadyExists", mrerr.ErrorKindUser, "rewrite name '{{ .name }}' already exists")

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrCompanyPageRewriteNameAlreadyExists,
	}
}
