package dictionaries

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	PaperColorAPIName = "Dictionaries.PaperColorAPI"
)

type (
	PaperColorAPI interface {
		// CheckingAvailability - error: FactoryErrPaperColorNotFound or Failed
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}
)
