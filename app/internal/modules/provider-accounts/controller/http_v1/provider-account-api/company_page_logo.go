package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/provider-account-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	companyPageItemLogoURL = "/v1/account/company-page/logo"
)

type (
	CompanyPageLogo struct {
		section mrcore.ClientSection
		service usecase.CompanyPageLogoService
	}
)

func NewCompanyPageLogo(
	section mrcore.ClientSection,
	service usecase.CompanyPageLogoService,
) *CompanyPageLogo {
	return &CompanyPageLogo{
		section: section,
		service: service,
	}
}

func (ht *CompanyPageLogo) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitCompanyPagePermission, next)
	}

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(companyPageItemLogoURL), moduleAccessFunc(ht.UploadLogo()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(companyPageItemLogoURL), moduleAccessFunc(ht.RemoveLogo()))
}

func (ht *CompanyPageLogo) UploadLogo() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		logger := mrctx.Logger(c.Context())

		file, hdr, err := c.Request().FormFile(module.ParamNameFileCompanyLogo)

		if err != nil {
			mrdebug.MultipartForm(logger, c.Request().MultipartForm)
			return mrcore.FactoryErrHttpMultipartFormFile.Wrap(err, module.ParamNameFileCompanyLogo)
		}

		defer file.Close()

		mrdebug.MultipartFileHeader(logger, hdr)

		item := mrtype.File{
			FileInfo: mrtype.FileInfo{
				ContentType:  hdr.Header.Get("Content-Type"),
				OriginalName: hdr.Filename,
				Size:         hdr.Size,
			},
			Body: file,
		}

		if err = ht.service.StoreFile(c.Context(), tmpAccountID, item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CompanyPageLogo) RemoveLogo() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.RemoveFile(c.Context(), tmpAccountID); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CompanyPageLogo) wrapError(err error, c mrcore.ClientContext) error {
	return err
}
