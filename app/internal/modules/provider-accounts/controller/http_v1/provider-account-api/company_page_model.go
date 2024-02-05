package http_v1

import entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

type (
	StoreCompanyPageRequest struct {
		RewriteName string `json:"rewriteName" validate:"required,max=64"`
		PageHead    string `json:"pageHead" validate:"required,max=128"`
		SiteURL     string `json:"siteUrl" validate:"required,max=256"`
	}

	ChangePublicStatusRequest struct {
		Status entity_shared.PublicStatus `json:"status" validate:"required"`
	}
)
