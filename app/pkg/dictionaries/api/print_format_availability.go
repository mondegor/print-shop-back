package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
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
		CheckingAvailability(ctx context.Context, itemID uint64) error
	}
)

var (
	// ErrPrintFormatRequired - print format ID is required.
	ErrPrintFormatRequired = mrerr.NewProto(
		"dictionaries.errPrintFormatRequired", mrerr.ErrorKindUser, "print format ID is required")

	// ErrPrintFormatNotAvailable - print format with ID is not available.
	ErrPrintFormatNotAvailable = mrerr.NewProto(
		"dictionaries.errPrintFormatNotAvailable", mrerr.ErrorKindUser, "print format with ID={{ .id }} is not available")

	// ErrPrintFormatNotFound - print format with ID not found.
	ErrPrintFormatNotFound = mrerr.NewProto(
		"dictionaries.errPrintFormatNotFound", mrerr.ErrorKindUser, "print format with ID={{ .id }} not found")
)
