package dictionaries

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	PaperColorAPIName = "Dictionaries.PaperColorAPI"
)

type (
	PaperColorAPI interface {
		// CheckingAvailability - error:
		//    - FactoryErrPaperColorRequired
		//	  - FactoryErrPaperColorNotAvailable
		//	  - FactoryErrPaperColorNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrPaperColorRequired = mrerr.NewFactory(
		"errDictionariesPaperColorRequired", mrerr.ErrorKindUser, "paper color ID is required")

	FactoryErrPaperColorNotAvailable = mrerr.NewFactory(
		"errDictionariesPaperColorNotAvailable", mrerr.ErrorKindUser, "paper color with ID={{ .id }} is not available")

	FactoryErrPaperColorNotFound = mrerr.NewFactory(
		"errDictionariesPaperColorNotFound", mrerr.ErrorKindUser, "paper color with ID={{ .id }} not found")
)
