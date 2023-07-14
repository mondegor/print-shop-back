package http_v1

import (
    "calc-user-data-back-adm/internal/controller/dto"
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/internal/usecase"
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrentity"
    "fmt"
    "net/http"
)

const (
    catalogBoxGetListURL = "/v1/catalog-boxes"
    catalogBoxGetItemURL = "/v1/catalog-boxes/:id"
    catalogBoxCreateURL = "/v1/catalog-boxes"
    catalogBoxStoreURL = "/v1/catalog-boxes/:id"
    catalogBoxChangeStatusURL = "/v1/catalog-boxes/:id/status"
    catalogBoxRemove = "/v1/catalog-boxes/:id"
)

type CatalogBox struct {
    service usecase.CatalogBoxService
}

func NewCatalogBox(service usecase.CatalogBoxService) *CatalogBox {
    return &CatalogBox{
        service: service,
    }
}

func (f *CatalogBox) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogBoxGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogBoxGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogBoxCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogBoxStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogBoxChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogBoxRemove, f.Remove())
}

func (f *CatalogBox) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *CatalogBox) newListFilter(c mrapp.ClientData) *entity.CatalogBoxListFilter {
    var listFilter entity.CatalogBoxListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *CatalogBox) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *CatalogBox) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogBox{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogBox{
            Article: request.Article,
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
            Depth: request.Depth,
        }

        err := f.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogBoxSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (f *CatalogBox) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogBox{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogBox{
            Id: f.getItemId(c),
            Version: request.Version,
            Article: request.Article,
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
            Depth: request.Depth,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogBox) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogBox{
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

func (f *CatalogBox) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogBox) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
