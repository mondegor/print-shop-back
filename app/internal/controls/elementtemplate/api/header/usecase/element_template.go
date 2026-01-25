package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
)

type (
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		storage      header.ElementTemplateStorage
		errorWrapper errors.Wrapper
		trace        mrtrace.Tracer
	}
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(
	storage header.ElementTemplateStorage,
	trace mrtrace.Tracer,
) *ElementTemplate {
	return &ElementTemplate{
		storage:      storage,
		errorWrapper: errors.NewUseCaseWrapper(),
		trace:        trace,
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
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
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
	uc.trace.Trace(
		ctx,
		"storage", api.ElementTemplateHeaderName,
		"cmd", command,
		"data", data,
	)
}
