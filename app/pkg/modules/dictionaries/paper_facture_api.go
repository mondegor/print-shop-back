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
		// CheckingAvailability - error:
		//    - FactoryErrPaperFactureRequired
		//	  - FactoryErrPaperFactureNotAvailable
		//	  - FactoryErrPaperFactureNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrPaperFactureRequired = mrerr.NewFactory(
		"errDictionariesPaperFactureRequired", mrerr.ErrorKindUser, "paper facture ID is required")

	FactoryErrPaperFactureNotAvailable = mrerr.NewFactory(
		"errDictionariesPaperFactureNotAvailable", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} is not available")

	FactoryErrPaperFactureNotFound = mrerr.NewFactory(
		"errDictionariesPaperFactureNotFound", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} not found")
)
