package header

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
)

type (
	// ElementTemplateStorage - comment interface.
	ElementTemplateStorage interface {
		FetchOneHead(ctx context.Context, rowID uint64) (entity.ElementTemplateHead, error)
	}
)
