package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

// ErrQueryHistoryNotFound - query with ShortLink not found.
var ErrQueryHistoryNotFound = mrerrfactory.NewProtoAppErrorByDefault(
	"errCalculationsQueryHistoryNotFound", mrerr.ErrorKindUser, "query {{ .shortLink }} not found")

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrQueryHistoryNotFound,
	}
}
