package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	// PaperColorAvailabilityName - название API.
	PaperColorAvailabilityName = "Dictionaries.API.PaperColorAvailability"
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
	ErrPaperColorRequired = mrerr.NewKindUser("PaperColorRequired", "paper color ID is required")

	// ErrPaperColorNotAvailable - paper color with ID is not available.
	ErrPaperColorNotAvailable = mrerr.NewKindUser("PaperColorNotAvailable", "paper color with ID={Id} is not available")

	// ErrPaperColorNotFound - paper color with ID not found.
	ErrPaperColorNotFound = mrerr.NewKindUser("PaperColorNotFound", "paper color with ID={Id} not found")
)
