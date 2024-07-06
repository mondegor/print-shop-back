package entity

import (
	"time"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"

	"github.com/google/uuid"
)

const (
	ModelNameCompanyPage     = "providers-api.ProviderAccounts.CompanyPage"     // ModelNameCompanyPage - название сущности
	ModelNameCompanyPageLogo = "providers-api.ProviderAccounts.CompanyPageLogo" // ModelNameCompanyPageLogo - название сущности
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct { // DB: printshop_providers.companies_pages
		AccountID   uuid.UUID         // account_id
		RewriteName string            `json:"rewriteName"`
		PageTitle   string            `json:"pageTitle"`
		LogoURL     string            `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL     string            `json:"siteUrl"`
		Status      enum.PublicStatus `json:"status"`
		CreatedAt   time.Time         `json:"createdAt"`
		UpdatedAt   *time.Time        `json:"updatedAt,omitempty"`
	}
)
