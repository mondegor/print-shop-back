package usecase_api

import (
	"context"
	entity "print-shop-back/internal/modules/controls/element-template/entity/admin-api"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplateStorage interface {
		FetchOneHead(ctx context.Context, rowID mrtype.KeyInt32) (entity.ElementTemplateHead, error)
	}
)
