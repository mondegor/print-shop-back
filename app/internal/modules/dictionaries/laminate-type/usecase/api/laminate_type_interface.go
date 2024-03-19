package usecase_api

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	LaminateTypeStorage interface {
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
