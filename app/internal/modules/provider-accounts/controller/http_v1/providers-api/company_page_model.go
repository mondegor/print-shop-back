package http_v1

import (
	"print-shop-back/pkg/modules/provider-accounts/enums"
)

type (
	StoreCompanyPageRequest struct {
		RewriteName string `json:"rewriteName" validate:"required,min=4,max=32,tag_rewrite_name"`
		PageTitle   string `json:"pageTitle" validate:"required,max=128"`
		SiteURL     string `json:"siteUrl" validate:"omitempty,max=512,http_url"`
	}

	ChangePublicStatusRequest struct {
		Status enums.PublicStatus `json:"status" validate:"required"`
	}
)
