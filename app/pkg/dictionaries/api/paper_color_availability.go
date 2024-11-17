package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
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
		CheckingAvailability(ctx context.Context, itemID uint64) error
	}
)

var (
	// ErrPaperColorRequired - paper color ID is required.
	ErrPaperColorRequired = mrerr.NewProto(
		"dictionaries.errPaperColorRequired", mrerr.ErrorKindUser, "paper color ID is required")

	// ErrPaperColorNotAvailable - paper color with ID is not available.
	ErrPaperColorNotAvailable = mrerr.NewProto(
		"dictionaries.errPaperColorNotAvailable", mrerr.ErrorKindUser, "paper color with ID={{ .id }} is not available")

	// ErrPaperColorNotFound - paper color with ID not found.
	ErrPaperColorNotFound = mrerr.NewProto(
		"dictionaries.errPaperColorNotFound", mrerr.ErrorKindUser, "paper color with ID={{ .id }} not found")
)
