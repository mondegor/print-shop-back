package entity

import (
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

const (
    ModelNameCompanyPageLogo = "CompanyPageLogo"
)

type (
    CompanyPageLogoObject struct {
        AccountId mrentity.KeyString
        File mrstorage.File
    }
)
