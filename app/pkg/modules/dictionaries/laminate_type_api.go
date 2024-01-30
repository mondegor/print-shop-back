package dictionaries

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	LaminateTypeAPIName = "Dictionaries.LaminateTypeAPI"
)

type (
	LaminateTypeAPI interface {
		// CheckingAvailability - error: FactoryErrLaminateTypeNotFound or Failed
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}
)
