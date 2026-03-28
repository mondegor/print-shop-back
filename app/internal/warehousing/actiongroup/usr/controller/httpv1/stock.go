package httpv1

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"github.com/mondegor/print-shop-back/internal/warehousing/module"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	stockListURL     = "/v1/warehousing/stocks"
	stockAddMoreURL  = "/v1/warehousing/stocks/add-more"
	stockMoveURL     = "/v1/warehousing/stocks/move"
	stockTransferURL = "/v1/warehousing/stocks/transfer"
)

type (
	// Stock - comment struct.
	Stock struct {
		parser               validate.RequestExtendParser
		sender               mrserver.ResponseSender
		serviceStock         usr.StockService
		useCaseMoveStock     moveStockUseCase
		useCaseAddMoreStock  addStockUseCase
		useCaseTransferStock transferStockUseCase
	}

	moveStockUseCase interface {
		Execute(ctx context.Context, item dto.MoveStockContainer) (movedStockID uint64, err error)
	}

	addStockUseCase interface {
		Execute(ctx context.Context, item dto.AddMoreStockContainer) (stockID uint64, err error)
	}

	transferStockUseCase interface {
		Execute(ctx context.Context, item dto.TransferStockContainers) error
	}
)

// NewStock - создаёт контроллер Stock.
func NewStock(
	parser validate.RequestExtendParser,
	sender mrserver.ResponseSender,
	serviceStock usr.StockService,
	useCaseMoveStock moveStockUseCase,
	useCaseAddStock addStockUseCase,
	useCaseTransferStock transferStockUseCase,
) *Stock {
	return &Stock{
		parser:               parser,
		sender:               sender,
		serviceStock:         serviceStock,
		useCaseMoveStock:     useCaseMoveStock,
		useCaseAddMoreStock:  useCaseAddStock,
		useCaseTransferStock: useCaseTransferStock,
	}
}

// Handlers - возвращает обработчики контроллера Stock.
func (ht *Stock) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: stockListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: stockAddMoreURL, Func: ht.AddMore},
		{Method: http.MethodPut, URL: stockMoveURL, Func: ht.Move},
		{Method: http.MethodPost, URL: stockTransferURL, Func: ht.Transfer},
	}
}

// GetList - comment method.
func (ht *Stock) GetList(w http.ResponseWriter, r *http.Request) error {
	items, hasNext, err := ht.serviceStock.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		StockListResponse{
			Items:   items,
			Cursor:  entity.CreateStockCursorValue(items),
			HasNext: hasNext,
		},
	)
}

func (ht *Stock) listParams(r *http.Request) dto.StockParams {
	return dto.StockParams{
		AccountID: ht.parser.UserID(r),
		Filter: dto.StockListFilter{
			SearchContainers: ht.parser.FilterUint64List(r, module.ParamNameFilterSearchStockContainers),
			SearchLocations:  ht.parser.FilterUint64List(r, module.ParamNameFilterSearchStockLocations),
		},
		Cursor: xtype.NewStockCursor(ht.parser.CursorParams(r)),
	}
}

// AddMore - comment method.
func (ht *Stock) AddMore(w http.ResponseWriter, r *http.Request) error {
	req := AddMoreStockRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := dto.AddMoreStockContainer{
		AccountID:   ht.parser.UserID(r),
		ContainerID: req.ContainerID,
		LocationID:  req.LocationID,
		Quantity:    req.Quantity,
	}

	itemID, err := ht.useCaseAddMoreStock.Execute(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		SuccessAddedStockResponse{
			StockID: itemID,
		},
	)
}

// Move - comment method.
func (ht *Stock) Move(w http.ResponseWriter, r *http.Request) error {
	req := MoveStockRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := dto.MoveStockContainer{
		AccountID:  ht.parser.UserID(r),
		StockID:    req.StockID,
		LocationID: req.LocationID,
		Quantity:   req.Quantity,
	}

	itemID, err := ht.useCaseMoveStock.Execute(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		SuccessMovedStockResponse{
			StockID: itemID,
		},
	)
}

// Transfer - comment method.
func (ht *Stock) Transfer(w http.ResponseWriter, r *http.Request) error {
	req := TransferStockRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := dto.TransferStockContainers{
		AccountID: ht.parser.UserID(r),
		Stocks: []dto.StockQuantity{
			{
				StockID:  req.StockID,
				Quantity: req.Quantity,
			},
		},
	}

	if err := ht.useCaseTransferStock.Execute(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Stock) getItemValue(r *http.Request) string {
	return ht.parser.PathParamString(r, "value")
}

func (ht *Stock) wrapError(err error, r *http.Request) error {
	if errors.Is(err, errors.ErrRecordNotFound) {
		return module.ErrStockNotFound.Wrap(err, ht.getItemValue(r))
	}

	return err
}
