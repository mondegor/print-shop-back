package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrentity"
)

const (
    catalogLaminateTypeGetListURL = "/v1/catalog-laminate-types"
    catalogLaminateTypeGetItemURL = "/v1/catalog-laminate-types/:id"
    catalogLaminateTypeCreateURL = "/v1/catalog-laminate-types"
    catalogLaminateTypeStoreURL = "/v1/catalog-laminate-types/:id"
    catalogLaminateTypeChangeStatusURL = "/v1/catalog-laminate-types/:id/status"
    catalogLaminateTypeRemove = "/v1/catalog-laminate-types/:id"
)

type CatalogLaminateType struct {
    service usecase.CatalogLaminateTypeService
}

func NewCatalogLaminateType(service usecase.CatalogLaminateTypeService) *CatalogLaminateType {
    return &CatalogLaminateType{
        service: service,
    }
}

func (f *CatalogLaminateType) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateTypeGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateTypeGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogLaminateTypeCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateTypeStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateTypeChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogLaminateTypeRemove, f.Remove())
}

func (f *CatalogLaminateType) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *CatalogLaminateType) newListFilter(c mrapp.ClientData) *entity.CatalogLaminateTypeListFilter {
    var listFilter entity.CatalogLaminateTypeListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *CatalogLaminateType) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *CatalogLaminateType) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogLaminateType{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminateType{
            Caption: request.Caption,
        }

        err := f.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogLaminateTypeSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (f *CatalogLaminateType) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogLaminateType{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminateType{
            Id: f.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogLaminateType) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminateType{
            Id: f.getItemId(c),
            Version: request.Version,
            Status: request.Status,
        }

        err := f.service.ChangeStatus(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogLaminateType) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogLaminateType) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
