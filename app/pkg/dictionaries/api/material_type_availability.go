package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
	"github.com/mondegor/go-webcore/mrtype"
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
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	// ErrMaterialTypeRequired - laminate type ID is required.
	ErrMaterialTypeRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesMaterialTypeRequired", mrerr.ErrorKindUser, "laminate type ID is required")

	// ErrMaterialTypeNotAvailable - laminate type with ID is not available.
	ErrMaterialTypeNotAvailable = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesMaterialTypeNotAvailable", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} is not available")

	// ErrMaterialTypeNotFound - laminate type with ID not found.
	ErrMaterialTypeNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errDictionariesMaterialTypeNotFound", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} not found")
)

// MaterialTypeErrors - comment func.
func MaterialTypeErrors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrMaterialTypeRequired,
		ErrMaterialTypeNotAvailable,
		ErrMaterialTypeNotFound,
	}
}
