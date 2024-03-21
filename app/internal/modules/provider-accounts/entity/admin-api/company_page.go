package entity

import (
	"print-shop-back/pkg/modules/provider-accounts/enums"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCompanyPage = "admin-api.ProviderAccounts.CompanyPage"
)

type (
	CompanyPage struct { // DB: printshop_providers.companies_pages
		AccountID   uuid.UUID          `json:"accountId"` // account_id
		RewriteName string             `json:"rewriteName" sort:"rewriteName"`
		PageTitle   string             `json:"pageTitle" sort:"pageTitle,default"`
		LogoURL     string             `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL     string             `json:"siteUrl" sort:"siteUrl"`
		Status      enums.PublicStatus `json:"status"`
		CreatedAt   time.Time          `json:"createdAt" sort:"createdAt"`
		UpdatedAt   *time.Time         `json:"updatedAt,omitempty" sort:"updatedAt"`
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
