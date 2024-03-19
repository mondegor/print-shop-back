package entity

import (
	"print-shop-back/pkg/modules/provider-accounts/enums"
	"time"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCompanyPage = "admin-api.ProviderAccounts.CompanyPage"
)

type (
	CompanyPage struct { // DB: printshop_provider_accounts.companies_pages
		AccountID mrtype.KeyString `json:"accountId"`                            // account_id
		UpdatedAt *time.Time       `json:"updatedAt,omitempty" sort:"updatedAt"` // updated_at

		RewriteName string `json:"rewriteName" sort:"rewriteName"`   // rewrite_name
		PageHead    string `json:"pageHead" sort:"pageHead,default"` // page_head
		LogoURL     string `json:"logoUrl,omitempty"`                // logo_meta.path
		SiteURL     string `json:"siteUrl" sort:"siteUrl"`           // site_url

		Status enums.PublicStatus `json:"status"` // page_status
	}

	CompanyPageParams struct {
		Filter CompanyPageListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	CompanyPageListFilter struct {
		SearchText string
		Statuses   []enums.PublicStatus
	}
)
