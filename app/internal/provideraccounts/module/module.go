package module

const (
	// Name - название модуля.
	Name = "ProviderAccounts"

	// Permission - разрешение модуля.
	Permission = "modProviderAccounts"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_providers"

	// DBTableNameCompaniesPages - таблица БД используемая модулем.
	DBTableNameCompaniesPages = DBSchema + ".companies_pages"

	// DBFieldWithoutDeletedAt - поле удаления записи не используется.
	DBFieldWithoutDeletedAt = ""

	// UnitCompanyPageLogoDir - относительный путь для хранения логотипов компаний.
	UnitCompanyPageLogoDir = "provider-account/companies-logo"
)
