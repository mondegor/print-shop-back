package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	// PrintFormatAvailabilityName - название API.
	PrintFormatAvailabilityName = "Dictionaries.API.PrintFormatAvailability"
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
	ErrPrintFormatRequired = mrerr.NewKindUser("PrintFormatRequired", "print format ID is required")

	// ErrPrintFormatNotAvailable - print format with ID is not available.
	ErrPrintFormatNotAvailable = mrerr.NewKindUser("PrintFormatNotAvailable", "print format with ID={Id} is not available")

	// ErrPrintFormatNotFound - print format with ID not found.
	ErrPrintFormatNotFound = mrerr.NewKindUser("PrintFormatNotFound", "print format with ID={Id} not found")
)
