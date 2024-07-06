package httpv1

import (
	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"
)

type (
	// StoreCompanyPageRequest - comment struct.
	StoreCompanyPageRequest struct {
		RewriteName string `json:"rewriteName" validate:"required,min=4,max=32,tag_rewrite_name"`
		PageTitle   string `json:"pageTitle" validate:"required,max=128"`
		SiteURL     string `json:"siteUrl" validate:"omitempty,max=512,http_url"`
	}

	// ChangePublicStatusRequest - comment struct.
	ChangePublicStatusRequest struct {
		Status enum.PublicStatus `json:"status" validate:"required"`
	}
)
