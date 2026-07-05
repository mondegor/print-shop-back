package entity

const (
	// ModelNameCompanyPage - название сущности.
	ModelNameCompanyPage = "public-api.ProviderAccounts.CompanyPage"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct { // DB: printshop_providers.companies_pages
		PageTitle string `json:"pageTitle"`
		LogoURL   string `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL   string `json:"siteUrl"`
	}
)
