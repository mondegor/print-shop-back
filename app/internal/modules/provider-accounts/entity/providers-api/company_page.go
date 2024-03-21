package entity

import (
	"print-shop-back/pkg/modules/provider-accounts/enums"
	"time"

	"github.com/google/uuid"
)

const (
	ModelNameCompanyPage     = "providers-api.ProviderAccounts.CompanyPage"
	ModelNameCompanyPageLogo = "providers-api.ProviderAccounts.CompanyPageLogo"
)

type (
	CompanyPage struct { // DB: printshop_providers.companies_pages
		AccountID   uuid.UUID          // account_id
		RewriteName string             `json:"rewriteName"`
		PageTitle   string             `json:"pageTitle"`
		LogoURL     string             `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL     string             `json:"siteUrl"`
		Status      enums.PublicStatus `json:"status"`
		CreatedAt   time.Time          `json:"createdAt"`
		UpdatedAt   *time.Time         `json:"updatedAt,omitempty"`
	}
)
