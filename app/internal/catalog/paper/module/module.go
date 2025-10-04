package module

const (
	// Name - название модуля.
	Name = "Catalog.Paper"

	// Permission - разрешение модуля.
	Permission = "modCatalogPaper"

	// LocaleDomain - домен локализации записей.
	LocaleDomain = "catalog.paper"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_catalog"

	// DBTableNamePapers - таблица БД используемая модулем.
	DBTableNamePapers = DBSchema + ".papers"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
