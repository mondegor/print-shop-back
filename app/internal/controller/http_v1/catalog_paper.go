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
    catalogPaperListURL = "/v1/catalog-papers"
    catalogPaperItemURL = "/v1/catalog-papers/:id"
    catalogPaperChangeStatusURL = "/v1/catalog-papers/:id/status"
)

type CatalogPaper struct {
    service usecase.CatalogPaperService
}

func NewCatalogPaper(service usecase.CatalogPaperService) *CatalogPaper {
    return &CatalogPaper{
        service: service,
    }
}

func (ht *CatalogPaper) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogPaperListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogPaperListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogPaperItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogPaperItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogPaperChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogPaper) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogPaper) newListFilter(c mrapp.ClientData) *entity.CatalogPaperListFilter {
    var listFilter entity.CatalogPaperListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogPaper) Get() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogPaper) Create() mrapp.HttpHandlerFunc {
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
            Sides: request.Sides,
        }

        err := ht.service.Create(c.Context(), &item)

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

func (ht *CatalogPaper) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogPaper{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
            Id:        ht.getItemId(c),
            Version:   request.Version,
            Article:   request.Article,
            Caption:   request.Caption,
            Length:    request.Length,
            Width:     request.Width,
            Density:   request.Density,
            ColorId:   request.ColorId,
            FactureId: request.FactureId,
            Thickness: request.Thickness,
            Sides:     request.Sides,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaper) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
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

func (ht *CatalogPaper) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaper) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
