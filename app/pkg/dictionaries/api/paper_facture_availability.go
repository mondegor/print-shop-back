package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	// PaperFactureAvailabilityName - название API.
	PaperFactureAvailabilityName = "Dictionaries.API.PaperFactureAvailability"
)

type (
	// PaperFactureAvailability - comment interface.
	PaperFactureAvailability interface {
		// CheckingAvailability - error:
		//    - ErrPaperFactureRequired
		//	  - ErrPaperFactureNotAvailable
		//	  - ErrPaperFactureNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID uint64) error
	}
)

var (
	// ErrPaperFactureRequired - paper facture ID is required.
	ErrPaperFactureRequired = mrerr.NewKindUser("PaperFactureRequired", "paper facture ID is required")

	// ErrPaperFactureNotAvailable - paper facture with ID is not available.
	ErrPaperFactureNotAvailable = mrerr.NewKindUser("PaperFactureNotAvailable", "paper facture with ID={Id} is not available")

	// ErrPaperFactureNotFound - paper facture with ID not found.
	ErrPaperFactureNotFound = mrerr.NewKindUser("PaperFactureNotFound", "paper facture with ID={Id} not found")
)
