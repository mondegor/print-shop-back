package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
)

type (
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		storage      header.ElementTemplateStorage
		errorWrapper mrerr.UseCaseErrorWrapper
		trace        mrtrace.Tracer
	}
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(
	storage header.ElementTemplateStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
	trace mrtrace.Tracer,
) *ElementTemplate {
	return &ElementTemplate{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, api.ElementTemplateHeaderName),
		trace:        trace,
	}
}

// GetItemHeader - comment method.
func (uc *ElementTemplate) GetItemHeader(ctx context.Context, itemID uint64) (api.ElementTemplateDTO, error) {
	uc.traceCmd(ctx, "GetHead", mrargs.Group{"id": itemID})

	if itemID == 0 {
		return api.ElementTemplateDTO{}, api.ErrElementTemplateRequired.New()
	}

	item, err := uc.storage.FetchOneHead(ctx, itemID)
	if err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return api.ElementTemplateDTO{}, api.ErrElementTemplateNotFound.New(itemID)
		}

		return api.ElementTemplateDTO{}, uc.errorWrapper.WrapErrorFailed(err)
	}

	if item.Status == mrenum.ItemStatusDisabled {
		return api.ElementTemplateDTO{}, api.ErrElementTemplateIsDisabled.New(itemID)
	}

	return item.ElementTemplateDTO, nil
}

func (uc *ElementTemplate) traceCmd(ctx context.Context, command string, data mrargs.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.ElementTemplateHeaderName,
		"cmd", command,
		"data", data,
	)
}
