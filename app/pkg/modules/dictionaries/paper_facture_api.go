package dictionaries

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	PaperFactureAPIName = "Dictionaries.PaperFactureAPI"
)

type (
	PaperFactureAPI interface {
		// CheckingAvailability - error: FactoryErrPaperFactureNotFound or Failed
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}
)
