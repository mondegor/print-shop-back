package module

const (
	Name       = "ProviderAccounts"    // Name - название модуля
	Permission = "modProviderAccounts" // Permission - разрешение модуля

	DBSchema                  = "printshop_providers" // DBSchema - схема БД используемая модулем
	DBTableNameCompaniesPages = "companies_pages"     // DBTableNameCompaniesPages - таблица БД используемая модулем

	UnitCompanyPagePermission = Permission                        // UnitCompanyPagePermission - разрешение юнита CompanyPage
	UnitCompanyPageLogoDir    = "provider-account/companies-logo" // UnitCompanyPageLogoDir - относительный путь для хранения логотипов компаний
)
