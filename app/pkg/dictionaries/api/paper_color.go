package api

import (
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/pkg/api"
)

const (
	// PaperColorAvailabilityName - название API.
	PaperColorAvailabilityName = "Dictionaries.API.PaperColorAvailability"
)

type (
	// PaperColorAvailability - проверяет доступность цвета бумаги по его ID.
	// CheckAvailability - error:
	//    - ErrPaperColorRequired
	//	  - ErrPaperColorNotAvailable
	//	  - ErrPaperColorNotFound
	//	  - Failed
	PaperColorAvailability api.AvailabilityChecker
)

var (
	// ErrPaperColorRequired - paper color ID is required.
	ErrPaperColorRequired = mrerr.NewKindUser("PaperColorRequired", "paper color ID is required")

	// ErrPaperColorNotAvailable - paper color with ID is not available.
	ErrPaperColorNotAvailable = mrerr.NewKindUser("PaperColorNotAvailable", "paper color with ID={Id} is not available")

	// ErrPaperColorNotFound - paper color with ID not found.
	ErrPaperColorNotFound = mrerr.NewKindUser("PaperColorNotFound", "paper color with ID={Id} not found")
)
