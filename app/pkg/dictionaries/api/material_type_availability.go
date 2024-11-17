package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	MaterialTypeAvailabilityName = "Dictionaries.API.MaterialTypeAvailability" // MaterialTypeAvailabilityName - название API
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
	ErrMaterialTypeRequired = mrerr.NewProto(
		"dictionaries.errMaterialTypeRequired", mrerr.ErrorKindUser, "laminate type ID is required")

	// ErrMaterialTypeNotAvailable - laminate type with ID is not available.
	ErrMaterialTypeNotAvailable = mrerr.NewProto(
		"dictionaries.errMaterialTypeNotAvailable", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} is not available")

	// ErrMaterialTypeNotFound - laminate type with ID not found.
	ErrMaterialTypeNotFound = mrerr.NewProto(
		"dictionaries.errMaterialTypeNotFound", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} not found")
)
