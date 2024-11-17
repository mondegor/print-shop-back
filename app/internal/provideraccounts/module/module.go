package module

const (
	Name       = "ProviderAccounts"    // Name - название модуля
	Permission = "modProviderAccounts" // Permission - разрешение модуля

	DBSchema                  = "printshop_providers"         // DBSchema - схема БД используемая модулем
	DBTableNameCompaniesPages = DBSchema + ".companies_pages" // DBTableNameCompaniesPages - таблица БД используемая модулем
	DBFieldWithoutDeletedAt   = ""                            // DBFieldWithoutDeletedAt - поле удаления записи не используется

	UnitCompanyPagePermission = Permission                        // UnitCompanyPagePermission - разрешение юнита CompanyPage
	UnitCompanyPageLogoDir    = "provider-account/companies-logo" // UnitCompanyPageLogoDir - относительный путь для хранения логотипов компаний
)
