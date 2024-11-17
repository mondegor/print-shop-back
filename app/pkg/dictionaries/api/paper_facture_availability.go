package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	PaperFactureAvailabilityName = "Dictionaries.API.PaperFactureAvailability" // PaperFactureAvailabilityName - название API
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
	ErrPaperFactureRequired = mrerr.NewProto(
		"dictionaries.errPaperFactureRequired", mrerr.ErrorKindUser, "paper facture ID is required")

	// ErrPaperFactureNotAvailable - paper facture with ID is not available.
	ErrPaperFactureNotAvailable = mrerr.NewProto(
		"dictionaries.errPaperFactureNotAvailable", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} is not available")

	// ErrPaperFactureNotFound - paper facture with ID not found.
	ErrPaperFactureNotFound = mrerr.NewProto(
		"dictionaries.errPaperFactureNotFound", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} not found")
)
