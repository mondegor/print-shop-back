package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

const (
    catalogLaminateTypeListURL = "/v1/catalog-laminate-types"
    catalogLaminateTypeItemURL = "/v1/catalog-laminate-types/:id"
    catalogLaminateTypeChangeStatusURL = "/v1/catalog-laminate-types/:id/status"
)

type CatalogLaminateType struct {
    service usecase.CatalogLaminateTypeService
}

func NewCatalogLaminateType(service usecase.CatalogLaminateTypeService) *CatalogLaminateType {
    return &CatalogLaminateType{
        service: service,
    }
}

func (ht *CatalogLaminateType) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateTypeListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogLaminateTypeListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogLaminateTypeItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateTypeItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogLaminateTypeItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogLaminateTypeChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogLaminateType) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogLaminateType) newListFilter(c mrcore.ClientData) *entity.CatalogLaminateTypeListFilter {
    var listFilter entity.CatalogLaminateTypeListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogLaminateType) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogLaminateType) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogLaminateType{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminateType{
            Caption: request.Caption,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := view.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogLaminateTypeSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogLaminateType) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogLaminateType{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminateType{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogLaminateType) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminateType{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Status:  request.Status,
        }

        err := ht.service.ChangeStatus(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogLaminateType) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogLaminateType) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
