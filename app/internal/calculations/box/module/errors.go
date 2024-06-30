package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

// ErrCalcResultNotFound - result with ID not found.
var ErrCalcResultNotFound = mrerrfactory.NewProtoAppErrorByDefault(
	"errCatalogCalcResultNotFound", mrerr.ErrorKindUser, "result with ID={{ .id }} not found")

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrCalcResultNotFound,
	}
}
