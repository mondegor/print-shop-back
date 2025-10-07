package api

import (
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/pkg/api"
)

const (
	// PrintFormatAvailabilityName - название API.
	PrintFormatAvailabilityName = "Dictionaries.API.PrintFormatAvailability"
)

type (
	// PrintFormatAvailability - проверяет доступность печатного формата по его ID.
	// CheckAvailability - error:
	//    - ErrPrintFormatRequired
	//	  - ErrPrintFormatNotAvailable
	//	  - ErrPrintFormatNotFound
	//	  - Failed
	PrintFormatAvailability api.AvailabilityChecker
)

var (
	// ErrPrintFormatRequired - print format ID is required.
	ErrPrintFormatRequired = mrerr.NewKindUser("PrintFormatRequired", "print format ID is required")

	// ErrPrintFormatNotAvailable - print format with ID is not available.
	ErrPrintFormatNotAvailable = mrerr.NewKindUser("PrintFormatNotAvailable", "print format with ID={Id} is not available")

	// ErrPrintFormatNotFound - print format with ID not found.
	ErrPrintFormatNotFound = mrerr.NewKindUser("PrintFormatNotFound", "print format with ID={Id} not found")
)
