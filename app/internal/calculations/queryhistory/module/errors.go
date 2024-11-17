package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

// ErrQueryHistoryNotFound - query with ShortLink not found.
var ErrQueryHistoryNotFound = mrerr.NewProto(
	"calculations.errQueryHistoryNotFound", mrerr.ErrorKindUser, "query {{ .shortLink }} not found")
