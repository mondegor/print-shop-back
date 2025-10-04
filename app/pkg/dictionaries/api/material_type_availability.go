package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	// MaterialTypeAvailabilityName - название API.
	MaterialTypeAvailabilityName = "Dictionaries.API.MaterialTypeAvailability"
)

type (
	// MaterialTypeAvailability - comment interface.
	MaterialTypeAvailability interface {
		// CheckingAvailability - error:
		//    - ErrMaterialTypeRequired
		//	  - ErrMaterialTypeNotAvailable
		//	  - ErrMaterialTypeNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID uint64) error
	}
)

var (
	// ErrMaterialTypeRequired - laminate type ID is required.
	ErrMaterialTypeRequired = mrerr.NewKindUser("MaterialTypeRequired", "laminate type ID is required")

	// ErrMaterialTypeNotAvailable - laminate type with ID is not available.
	ErrMaterialTypeNotAvailable = mrerr.NewKindUser("MaterialTypeNotAvailable", "laminate type with ID={Id} is not available")

	// ErrMaterialTypeNotFound - laminate type with ID not found.
	ErrMaterialTypeNotFound = mrerr.NewKindUser("MaterialTypeNotFound", "laminate type with ID={Id} not found")
)
