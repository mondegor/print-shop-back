package usecase

import (
	"context"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrworkflow/itemstatus"
	"github.com/mondegor/go-core/util/conv"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/controls/elementtemplate/api/header"
	"print-shop-back/pkg/controls/api"
)

type (
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		storage      header.ElementTemplateStorage
		errorWrapper errors.Wrapper
		tracer       trace.Tracer
	}
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(
	storage header.ElementTemplateStorage,
	tracer trace.Tracer,
) *ElementTemplate {
	return &ElementTemplate{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
		tracer:       tracer,
	}
}

// GetItemHeader - comment method.
func (uc *ElementTemplate) GetItemHeader(ctx context.Context, itemID uint64) (api.ElementTemplateDTO, error) {
	uc.traceCmd(ctx, "GetHead", conv.Group{"id": itemID})

	if itemID == 0 {
		return api.ElementTemplateDTO{}, api.ErrElementTemplateRequired
	}

	item, err := uc.storage.FetchOneHead(ctx, itemID)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return api.ElementTemplateDTO{}, api.ErrElementTemplateNotFound.Wrap(err, itemID)
		}

		return api.ElementTemplateDTO{}, uc.errorWrapper.Wrap(err)
	}

	if item.Status == itemstatus.Disabled {
		return api.ElementTemplateDTO{}, api.ErrElementTemplateIsDisabled.New(itemID)
	}

	return item.ElementTemplateDTO, nil
}

func (uc *ElementTemplate) traceCmd(ctx context.Context, command string, data conv.Group) {
	uc.tracer.Trace(
		ctx,
		"storage", api.ElementTemplateHeaderName,
		"cmd", command,
		"data", data,
	)
}
