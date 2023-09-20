package http_v1

import (
    "net/http"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

const (
    tmpAccountId = "b8a75cd5-ddd1-46ab-be84-406e669cbfa9"

    accountCompanyPageItemURL = "/v1/account/company-page"
    accountCompanyPageChangeStatusURL = "/v1/account/company-page/status"
    accountCompanyPageItemLogoURL = "/v1/account/company-page/logo"
)

type (
    AccountCompanyPage struct {
        service usecase.AccountCompanyPageService
        serviceLogo usecase.AccountCompanyPageLogoService
    }
)

func NewAccountCompanyPage(service usecase.AccountCompanyPageService,
                           serviceLogo usecase.AccountCompanyPageLogoService) *AccountCompanyPage {
    return &AccountCompanyPage{
        service: service,
        serviceLogo: serviceLogo,
    }
}

func (ht *AccountCompanyPage) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, accountCompanyPageItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, accountCompanyPageItemURL, ht.Store())

    router.HttpHandlerFunc(http.MethodPut, accountCompanyPageChangeStatusURL, ht.ChangeStatus())

    router.HttpHandlerFunc(http.MethodPut, accountCompanyPageItemLogoURL, ht.UploadLogo())
    router.HttpHandlerFunc(http.MethodDelete, accountCompanyPageItemLogoURL, ht.RemoveLogo())
}

func (ht *AccountCompanyPage) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), tmpAccountId)

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *AccountCompanyPage) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreAccountCompanyPageRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CompanyPage{
            AccountId: tmpAccountId,
            Version: request.Version,
            RewriteName: request.RewriteName,
            PageHead: request.PageHead,
            SiteUrl: request.SiteUrl,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *AccountCompanyPage) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeResourceStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CompanyPage{
            AccountId: tmpAccountId,
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

func (ht *AccountCompanyPage) UploadLogo() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        file, hdr, err := c.Request().FormFile("companyLogo")

        if err != nil {
            return err
        }

        defer file.Close()

        logger := mrctx.Logger(c.Context())

        logger.Debug(
            "uploaded file: name=%s, size=%d, header=%#v",
            hdr.Filename, hdr.Size, hdr.Header,
        )

        item := entity.CompanyPageLogoObject{
            AccountId: tmpAccountId,
            File: mrstorage.File{
                ContentType: hdr.Header.Get("Content-Type"),
                Name: hdr.Filename,
                Size: hdr.Size,
                Body: file,
            },
        }

        err = ht.serviceLogo.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *AccountCompanyPage) RemoveLogo() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.serviceLogo.Remove(c.Context(), tmpAccountId)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}
