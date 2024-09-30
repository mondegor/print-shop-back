package header

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
)

type (
	// ElementTemplateStorage - comment interface.
	ElementTemplateStorage interface {
		FetchOneHead(ctx context.Context, rowID mrtype.KeyInt32) (entity.ElementTemplateHead, error)
	}
)
