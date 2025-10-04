package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"
)

const (
	// ModelNameCompanyPage - название сущности.
	ModelNameCompanyPage = "admin-api.ProviderAccounts.CompanyPage"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct { // DB: printshop_providers.companies_pages
		AccountID   uuid.UUID         `json:"accountId"` // account_id
		RewriteName string            `json:"rewriteName" sort:"rewriteName"`
		PageTitle   string            `json:"pageTitle" sort:"pageTitle,default"`
		LogoURL     string            `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL     string            `json:"siteUrl" sort:"siteUrl"`
		Status      enum.PublicStatus `json:"status"`
		CreatedAt   time.Time         `json:"createdAt" sort:"createdAt"`
		UpdatedAt   time.Time         `json:"updatedAt" sort:"updatedAt"`
	}

	// CompanyPageParams - comment struct.
	CompanyPageParams struct {
		Filter CompanyPageListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// CompanyPageListFilter - comment struct.
	CompanyPageListFilter struct {
		SearchText string
		Statuses   []enum.PublicStatus
	}
)
