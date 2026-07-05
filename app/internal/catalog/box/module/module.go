package module

const (
	// Name - название модуля.
	Name = "Catalog.Box"

	// Permission - разрешение модуля.
	Permission = "modCatalogBox"

	// LocaleDomain - домен локализации записей.
	LocaleDomain = "catalog.box"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_catalog"

	// DBTableNameBoxes - таблица БД используемая модулем.
	DBTableNameBoxes = DBSchema + ".boxes"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
