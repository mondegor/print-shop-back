package api

import (
	"github.com/mondegor/go-sysmess/errors"

	"print-shop-back/pkg/api"
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
	ErrPaperColorRequired = errors.NewUserProto("PaperColorRequired", "paper color ID is required")

	// ErrPaperColorNotAvailable - paper color with ID is not available.
	ErrPaperColorNotAvailable = errors.NewUserProto("PaperColorNotAvailable", "paper color with ID={Id} is not available")

	// ErrPaperColorNotFound - paper color with ID not found.
	ErrPaperColorNotFound = errors.NewUserProto("PaperColorNotFound", "paper color with ID={Id} not found")
)
