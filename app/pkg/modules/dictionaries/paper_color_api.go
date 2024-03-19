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
		// CheckingAvailability - error: FactoryErrPaperColorRequired | FactoryErrPaperColorNotFound | Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrPaperColorRequired = mrerr.NewFactory(
		"errDictionariesPaperColorRequired", mrerr.ErrorKindUser, "paper color ID is required")

	FactoryErrPaperColorNotFound = mrerr.NewFactory(
		"errDictionariesPaperColorNotFound", mrerr.ErrorKindUser, "paper color with ID={{ .id }} not found")
)
