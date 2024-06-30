package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// PaperFactureStorage - comment interface.
	PaperFactureStorage interface {
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
	}
)
