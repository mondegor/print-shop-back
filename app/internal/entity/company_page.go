package entity

import (
    "time"

    "github.com/mondegor/go-storage/mrentity"
)

const (
    ModelNameCompanyPage = "CompanyPage"
)

type (
    CompanyPage struct { // DB: accounts_companies_pages
        AccountId   mrentity.KeyString `json:"accountId"` // account_id
        Version     mrentity.Version `json:"version"` // tag_version
        UpdateAt    time.Time `json:"updateAt"` // datetime_updated
        RewriteName string `json:"rewriteName"` // rewrite_name
        PageHead    string `json:"pageHead"` // page_head
        LogoPath    string `json:"logoPath"` // logo_path
        SiteUrl     string `json:"siteUrl"` // site_url
        Status      ResourceStatus `json:"status"` // page_status
    }

    CompanyPageListFilter struct {
        Statuses  []ResourceStatus
    }
)
