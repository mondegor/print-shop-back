package entity

const (
	ModelNameCompanyPage = "public-api.ProviderAccounts.CompanyPage"
)

type (
	CompanyPage struct { // DB: ps_provider_accounts.companies_pages
		PageHead string `json:"pageHead"`          // page_head
		LogoURL  string `json:"logoUrl,omitempty"` // logo_meta.path
		SiteURL  string `json:"siteUrl"`           // site_url
	}
)
