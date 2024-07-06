package entity

const (
	ModelNameCompanyPage = "public-api.ProviderAccounts.CompanyPage" // ModelNameCompanyPage - название сущности
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct { // DB: printshop_providers.companies_pages
		PageTitle string `json:"pageTitle"`
		LogoURL   string `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL   string `json:"siteUrl"`
	}
)
