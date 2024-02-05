package entity

import (
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"
	"time"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCompanyPage     = "provider-account-api.ProviderAccounts.CompanyPage"
	ModelNameCompanyPageLogo = "provider-account-api.ProviderAccounts.CompanyPageLogo"
)

type (
	CompanyPage struct { // DB: printdata_provider_accounts.companies_pages
		AccountID mrtype.KeyString // account_id
		UpdatedAt time.Time        `json:"updatedAt"` // updated_at

		RewriteName string `json:"rewriteName"`       // rewrite_name
		PageHead    string `json:"pageHead"`          // page_head
		LogoURL     string `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL     string `json:"siteUrl"`           // site_url

		Status entity_shared.PublicStatus `json:"status"` // page_status
	}
)
