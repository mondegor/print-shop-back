package availability

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// MaterialTypeStorage - comment interface.
	MaterialTypeStorage interface {
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
	}
)
