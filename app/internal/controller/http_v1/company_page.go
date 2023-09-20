package http_v1

import (
    "net/http"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    companyPageListURL = "/v1/accounts/companies-pages"
)

type (
    CompanyPage struct {
        service usecase.CompanyPageService
    }
)

func NewCompanyPage(service usecase.CompanyPageService) *CompanyPage {
    return &CompanyPage{
        service: service,
    }
}

func (ht *CompanyPage) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, companyPageListURL, ht.GetList())
}

func (ht *CompanyPage) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CompanyPage) newListFilter(c mrcore.ClientData) *entity.CompanyPageListFilter {
    var listFilter entity.CompanyPageListFilter

    parseFilterResourceStatuses(c, &listFilter.Statuses)

    return &listFilter
}

//func (ht *CompanyPage) GetLogo() mrcore.HttpHandlerFunc {
//    return func(c mrcore.ClientData) error {
//        item, err := ht.serviceLogo.GetItem(c.Context(), tmpAccountId)
//
//        if err != nil {
//            return err
//        }
//
//        defer item.File.Body.Close()
//
//        return c.SendFile(item.File.ContentType, item.File.Body)
//    }
//}
