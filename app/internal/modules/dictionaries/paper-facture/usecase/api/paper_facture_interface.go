package usecase_api

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperFactureStorage interface {
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
