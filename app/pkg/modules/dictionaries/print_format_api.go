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
		// CheckingAvailability - error: FactoryErrPrintFormatRequired | FactoryErrPrintFormatNotFound | Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrPrintFormatRequired = mrerr.NewFactory(
		"errDictionariesPrintFormatRequired", mrerr.ErrorKindUser, "print format ID is required")

	FactoryErrPrintFormatNotFound = mrerr.NewFactory(
		"errDictionariesPrintFormatNotFound", mrerr.ErrorKindUser, "print format with ID={{ .id }} not found")
)
