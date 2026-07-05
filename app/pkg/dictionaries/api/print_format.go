package api

import (
	"github.com/mondegor/go-sysmess/errors"

	"print-shop-back/pkg/api"
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
	ErrPrintFormatRequired = errors.NewUserProto("PrintFormatRequired", "print format ID is required")

	// ErrPrintFormatNotAvailable - print format with ID is not available.
	ErrPrintFormatNotAvailable = errors.NewUserProto("PrintFormatNotAvailable", "print format with ID={Id} is not available")

	// ErrPrintFormatNotFound - print format with ID not found.
	ErrPrintFormatNotFound = errors.NewUserProto("PrintFormatNotFound", "print format with ID={Id} not found")
)
