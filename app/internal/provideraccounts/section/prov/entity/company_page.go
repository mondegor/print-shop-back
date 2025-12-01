package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/type/publicstatus"
)

const (
	// ModelNameCompanyPage - название сущности.
	ModelNameCompanyPage = "providers-api.ProviderAccounts.CompanyPage"

	// ModelNameCompanyPageLogo - название сущности.
	ModelNameCompanyPageLogo = "providers-api.ProviderAccounts.CompanyPageLogo"
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
