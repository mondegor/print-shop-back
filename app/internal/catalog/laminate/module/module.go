package module

const (
	// Name - название модуля.
	Name = "Catalog.Laminate"

	// Permission - разрешение модуля.
	Permission = "modCatalogLaminate"

	// LocaleDomain - домен локализации записей.
	LocaleDomain = "catalog.laminate"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_catalog"

	// DBTableNameLaminates - таблица БД используемая модулем.
	DBTableNameLaminates = DBSchema + ".laminates"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
