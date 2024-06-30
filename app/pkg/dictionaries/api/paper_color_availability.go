package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	PaperColorAvailabilityName = "Dictionaries.API.PaperColorAvailability" // PaperColorAvailabilityName - название API
)

type (
	// PaperColorAvailability - comment interface.
	PaperColorAvailability interface {
		// CheckingAvailability - error:
		//    - ErrPaperColorRequired
		//	  - ErrPaperColorNotAvailable
		//	  - ErrPaperColorNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	// ErrPaperColorRequired - paper color ID is required.
	ErrPaperColorRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPaperColorRequired", mrerr.ErrorKindUser, "paper color ID is required")

	// ErrPaperColorNotAvailable - paper color with ID is not available.
	ErrPaperColorNotAvailable = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPaperColorNotAvailable", mrerr.ErrorKindUser, "paper color with ID={{ .id }} is not available")

	// ErrPaperColorNotFound - paper color with ID not found.
	ErrPaperColorNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesPaperColorNotFound", mrerr.ErrorKindUser, "paper color with ID={{ .id }} not found")
)

// PaperColorErrors - comment func.
func PaperColorErrors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrPaperColorRequired,
		ErrPaperColorNotAvailable,
		ErrPaperColorNotFound,
	}
}
