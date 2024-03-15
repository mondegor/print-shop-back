package dictionaries

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	PrintFormatAPIName = "Dictionaries.PrintFormatAPI"
)

type (
	PrintFormatAPI interface {
		// CheckingAvailability - error: FactoryErrPrintFormatNotFound or Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrPrintFormatNotFound = mrerr.NewFactory(
		"errDictionariesPrintFormatNotFound", mrerr.ErrorKindUser, "print format with ID={{ .id }} not found")
)
