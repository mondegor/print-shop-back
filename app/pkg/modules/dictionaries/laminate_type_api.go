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
		// CheckingAvailability - error: FactoryErrLaminateTypeNotFound or Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrLaminateTypeNotFound = mrerr.NewFactory(
		"errDictionariesLaminateTypeNotFound", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} not found")
)
