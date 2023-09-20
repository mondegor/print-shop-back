package view

import (
    "github.com/mondegor/go-storage/mrentity"
)

type (
    StoreAccountCompanyPageRequest struct {
        Version  mrentity.Version `json:"version" validate:"required,gte=1"`
        RewriteName string `json:"rewriteName" validate:"required,max=64"`
        PageHead    string `json:"pageHead" validate:"required,max=128"`
        SiteUrl     string `json:"siteUrl" validate:"required,max=256"`
    }

    PublicCompanyPageResponse struct {
        PageHead string `json:"pageHead"`
        LogoPath string `json:"logoPath"`
        SiteUrl  string `json:"siteUrl"`
    }
)
