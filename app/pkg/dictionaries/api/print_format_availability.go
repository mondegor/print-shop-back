package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	PrintFormatAvailabilityName = "Dictionaries.API.PrintFormatAvailability" // PrintFormatAvailabilityName - название API
)

type (
	// PrintFormatAvailability - comment interface.
	PrintFormatAvailability interface {
		// CheckingAvailability - error:
		//    - ErrPrintFormatRequired
		//	  - ErrPrintFormatNotAvailable
		//	  - ErrPrintFormatNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	// ErrPrintFormatRequired - print format ID is required.
	ErrPrintFormatRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPrintFormatRequired", mrerr.ErrorKindUser, "print format ID is required")

	// ErrPrintFormatNotAvailable - print format with ID is not available.
	ErrPrintFormatNotAvailable = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPrintFormatNotAvailable", mrerr.ErrorKindUser, "print format with ID={{ .id }} is not available")

	// ErrPrintFormatNotFound - print format with ID not found.
	ErrPrintFormatNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPrintFormatNotFound", mrerr.ErrorKindUser, "print format with ID={{ .id }} not found")
)

// PrintFormatErrors - comment func.
func PrintFormatErrors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrPrintFormatRequired,
		ErrPrintFormatNotAvailable,
		ErrPrintFormatNotFound,
	}
}
