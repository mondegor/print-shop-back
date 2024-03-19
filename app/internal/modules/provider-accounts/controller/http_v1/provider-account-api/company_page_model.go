package http_v1

import (
	"print-shop-back/pkg/modules/provider-accounts/enums"
)

type (
	StoreCompanyPageRequest struct {
		RewriteName string `json:"rewriteName" validate:"required,max=64"`
		PageHead    string `json:"pageHead" validate:"required,max=128"`
		SiteURL     string `json:"siteUrl" validate:"required,max=256"`
	}

	ChangePublicStatusRequest struct {
		Status enums.PublicStatus `json:"status" validate:"required"`
	}
)
