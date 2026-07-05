package api

import (
	"github.com/mondegor/go-core/errors"

	"print-shop-back/pkg/api"
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
	ErrPaperFactureRequired = errors.NewUserProto("PaperFactureRequired", "paper facture ID is required")

	// ErrPaperFactureNotAvailable - paper facture with ID is not available.
	ErrPaperFactureNotAvailable = errors.NewUserProto("PaperFactureNotAvailable", "paper facture with ID={Id} is not available")

	// ErrPaperFactureNotFound - paper facture with ID not found.
	ErrPaperFactureNotFound = errors.NewUserProto("PaperFactureNotFound", "paper facture with ID={Id} not found")
)
