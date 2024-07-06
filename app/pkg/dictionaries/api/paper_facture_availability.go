package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
	"github.com/mondegor/go-webcore/mrtype"
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
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	// ErrPaperFactureRequired - paper facture ID is required.
	ErrPaperFactureRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPaperFactureRequired", mrerr.ErrorKindUser, "paper facture ID is required")

	// ErrPaperFactureNotAvailable - paper facture with ID is not available.
	ErrPaperFactureNotAvailable = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPaperFactureNotAvailable", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} is not available")

	// ErrPaperFactureNotFound - paper facture with ID not found.
	ErrPaperFactureNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPaperFactureNotFound", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} not found")
)

// PaperFactureErrors - comment func.
func PaperFactureErrors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrPaperFactureRequired,
		ErrPaperFactureNotAvailable,
		ErrPaperFactureNotFound,
	}
}
