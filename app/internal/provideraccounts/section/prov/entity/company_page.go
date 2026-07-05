package entity

import (
	"time"

	"github.com/google/uuid"

	"print-shop-back/pkg/provideraccounts/enum/publicstatus"
)

const (
	// ModelNameCompanyPage - название сущности.
	ModelNameCompanyPage = "provider-api.ProviderAccounts.CompanyPage"

	// ModelNameCompanyPageLogo - название сущности.
	ModelNameCompanyPageLogo = "provider-api.ProviderAccounts.CompanyPageLogo"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct { // DB: printshop_providers.companies_pages
		AccountID   uuid.UUID         `json:"-"` // account_id
		RewriteName string            `json:"rewriteName"`
		PageTitle   string            `json:"pageTitle"`
		LogoURL     string            `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL     string            `json:"siteUrl"`
		Status      publicstatus.Enum `json:"status"`
		CreatedAt   time.Time         `json:"createdAt"`
		UpdatedAt   time.Time         `json:"updatedAt"`
	}
)
