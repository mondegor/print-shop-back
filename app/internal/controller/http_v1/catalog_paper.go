package http_v1

import (
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrentity"
    "fmt"
    "net/http"
)

const (
    catalogPaperGetListURL = "/v1/catalog-papers"
    catalogPaperGetItemURL = "/v1/catalog-papers/:id"
    catalogPaperCreateURL = "/v1/catalog-papers"
    catalogPaperStoreURL = "/v1/catalog-papers/:id"
    catalogPaperChangeStatusURL = "/v1/catalog-papers/:id/status"
    catalogPaperRemove = "/v1/catalog-papers/:id"
)

type CatalogPaper struct {
    service usecase.CatalogPaperService
}

func NewCatalogPaper(service usecase.CatalogPaperService) *CatalogPaper {
    return &CatalogPaper{
        service: service,
    }
}

func (f *CatalogPaper) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogPaperGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogPaperGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogPaperCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogPaperRemove, f.Remove())
}

func (f *CatalogPaper) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *CatalogPaper) newListFilter(c mrapp.ClientData) *entity.CatalogPaperListFilter {
    var listFilter entity.CatalogPaperListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *CatalogPaper) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *CatalogPaper) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogPaper{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
            Article: request.Article,
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
            Density: request.Density,
            ColorId: request.ColorId,
            FactureId: request.FactureId,
            Thickness: request.Thickness,
        }

        err := f.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogPaperSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (f *CatalogPaper) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogPaper{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
            Id: f.getItemId(c),
            Version: request.Version,
            Article: request.Article,
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
            Density: request.Density,
            ColorId: request.ColorId,
            FactureId: request.FactureId,
            Thickness: request.Thickness,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPaper) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
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

func (f *CatalogPaper) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPaper) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
