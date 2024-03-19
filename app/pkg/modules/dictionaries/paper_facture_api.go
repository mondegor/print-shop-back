package dictionaries

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	PaperFactureAPIName = "Dictionaries.PaperFactureAPI"
)

type (
	PaperFactureAPI interface {
		// CheckingAvailability - error: FactoryErrPaperFactureRequired | FactoryErrPaperFactureNotFound | Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrPaperFactureRequired = mrerr.NewFactory(
		"errDictionariesPaperFactureRequired", mrerr.ErrorKindUser, "paper facture ID is required")

	FactoryErrPaperFactureNotFound = mrerr.NewFactory(
		"errDictionariesPaperFactureNotFound", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} not found")
)
