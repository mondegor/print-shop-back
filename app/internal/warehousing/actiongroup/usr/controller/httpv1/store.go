package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"github.com/mondegor/print-shop-back/internal/warehousing/module"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	storeListURL = "/v1/warehousing/stores"
)

type (
	// Store - comment struct.
	Store struct {
		parser       validate.RequestExtendParser
		sender       mrserver.ResponseSender
		serviceStore usr.StoreService
	}
)

// NewStore - создаёт контроллер Store.
func NewStore(
	parser validate.RequestExtendParser,
	sender mrserver.ResponseSender,
	serviceStore usr.StoreService,
) *Store {
	return &Store{
		parser:       parser,
		sender:       sender,
		serviceStore: serviceStore,
	}
}

// Handlers - возвращает обработчики контроллера Store.
func (ht *Store) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: storeListURL, Func: ht.GetList},
	}
}

// GetList - comment method.
func (ht *Store) GetList(w http.ResponseWriter, r *http.Request) error {
	items, hasNext, err := ht.serviceStore.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		StoreListResponse{
			Items:   items,
			Cursor:  entity.CreateStoreCursorValue(items),
			HasNext: hasNext,
		},
	)
}

func (ht *Store) listParams(r *http.Request) dto.StoreParams {
	return dto.StoreParams{
		AccountID: ht.parser.UserID(r),
		Filter: dto.StoreListFilter{
			SearchCode:        ht.parser.FilterString(r, module.ParamNameFilterSearchStoreCode),
			SearchTerritories: ht.parser.FilterUint64List(r, module.ParamNameFilterSearchTerritories),
		},
		Cursor: xtype.NewStoreCursor(ht.parser.CursorParams(r)),
	}
}
