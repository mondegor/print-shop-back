package http_v1

import (
    "net/http"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    publicCompanyPageItemURL = "/v1/company/:rewrite-name"
)

type (
    PublicCompanyPage struct {
        service usecase.PublicCompanyPageService
    }
)

func NewPublicCompanyPage(service usecase.PublicCompanyPageService) *PublicCompanyPage {
    return &PublicCompanyPage{
        service: service,
    }
}

func (ht *PublicCompanyPage) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, publicCompanyPageItemURL, ht.Get())
}

func (ht *PublicCompanyPage) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItemByName(c.Context(), c.RequestPath().Get("rewrite-name"))

        if err != nil {
            return err
        }

        return c.SendResponse(
            http.StatusOK,
            view.PublicCompanyPageResponse{
                PageHead: item.PageHead,
                LogoPath: item.LogoPath,
                SiteUrl: item.SiteUrl,
            },
        )
    }
}
