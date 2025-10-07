package api

import (
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/pkg/api"
)

const (
	// MaterialTypeAvailabilityName - название API.
	MaterialTypeAvailabilityName = "Dictionaries.API.MaterialTypeAvailability"
)

type (
	// MaterialTypeAvailability - проверяет доступность типа материала по его ID.
	// CheckAvailability - error:
	//    - ErrMaterialTypeRequired
	//	  - ErrMaterialTypeNotAvailable
	//	  - ErrMaterialTypeNotFound
	//	  - Failed
	MaterialTypeAvailability api.AvailabilityChecker
)

var (
	// ErrMaterialTypeRequired - laminate type ID is required.
	ErrMaterialTypeRequired = mrerr.NewKindUser("MaterialTypeRequired", "laminate type ID is required")

	// ErrMaterialTypeNotAvailable - laminate type with ID is not available.
	ErrMaterialTypeNotAvailable = mrerr.NewKindUser("MaterialTypeNotAvailable", "laminate type with ID={Id} is not available")

	// ErrMaterialTypeNotFound - laminate type with ID not found.
	ErrMaterialTypeNotFound = mrerr.NewKindUser("MaterialTypeNotFound", "laminate type with ID={Id} not found")
)
