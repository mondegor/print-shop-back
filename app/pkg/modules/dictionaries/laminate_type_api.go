package dictionaries

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	LaminateTypeAPIName = "Dictionaries.LaminateTypeAPI"
)

type (
	LaminateTypeAPI interface {
		// CheckingAvailability - error:
		//    - FactoryErrLaminateTypeRequired
		//	  - FactoryErrLaminateTypeNotAvailable
		//	  - FactoryErrLaminateTypeNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrLaminateTypeRequired = mrerr.NewFactory(
		"errDictionariesLaminateTypeRequired", mrerr.ErrorKindUser, "laminate type ID is required")

	FactoryErrLaminateTypeNotAvailable = mrerr.NewFactory(
		"errDictionariesLaminateTypeNotAvailable", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} is not available")

	FactoryErrLaminateTypeNotFound = mrerr.NewFactory(
		"errDictionariesLaminateTypeNotFound", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} not found")
)
