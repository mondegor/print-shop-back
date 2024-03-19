package entity

import (
	"print-shop-back/pkg/modules/provider-accounts/enums"
	"time"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCompanyPage     = "provider-account-api.ProviderAccounts.CompanyPage"
	ModelNameCompanyPageLogo = "provider-account-api.ProviderAccounts.CompanyPageLogo"
)

type (
	CompanyPage struct { // DB: printshop_provider_accounts.companies_pages
		AccountID mrtype.KeyString // account_id
		UpdatedAt time.Time        `json:"updatedAt"` // updated_at

		RewriteName string `json:"rewriteName"`       // rewrite_name
		PageHead    string `json:"pageHead"`          // page_head
		LogoURL     string `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL     string `json:"siteUrl"`           // site_url

		Status enums.PublicStatus `json:"status"` // page_status
	}
)
