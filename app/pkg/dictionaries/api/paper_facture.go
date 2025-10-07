package api

import (
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/pkg/api"
)

const (
	// PaperFactureAvailabilityName - название API.
	PaperFactureAvailabilityName = "Dictionaries.API.PaperFactureAvailability"
)

type (
	// PaperFactureAvailability - проверяет доступность фактуры бумаги по его ID.
	// CheckAvailability - error:
	//    - ErrPaperFactureRequired
	//	  - ErrPaperFactureNotAvailable
	//	  - ErrPaperFactureNotFound
	//	  - Failed
	PaperFactureAvailability api.AvailabilityChecker
)

var (
	// ErrPaperFactureRequired - paper facture ID is required.
	ErrPaperFactureRequired = mrerr.NewKindUser("PaperFactureRequired", "paper facture ID is required")

	// ErrPaperFactureNotAvailable - paper facture with ID is not available.
	ErrPaperFactureNotAvailable = mrerr.NewKindUser("PaperFactureNotAvailable", "paper facture with ID={Id} is not available")

	// ErrPaperFactureNotFound - paper facture with ID not found.
	ErrPaperFactureNotFound = mrerr.NewKindUser("PaperFactureNotFound", "paper facture with ID={Id} not found")
)
