package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
// // ErrBoxNotFound - box with ID not found.
// ErrBoxNotFound = mrerrfactory.NewProtoAppErrorByDefault(
//
//	"errCalculationsBoxNotFound", mrerr.ErrorKindUser, "box with ID={{ .id }} not found")
)

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return nil
}
