package module

const (
	Name       = "Catalog.Box"   // Name - название модуля
	Permission = "modCatalogBox" // Permission - разрешение модуля

	DBSchema          = "printshop_catalog" // DBSchema - схема БД используемая модулем
	DBTableNameBoxes  = DBSchema + ".boxes" // DBTableNameBoxes - таблица БД используемая модулем
	DBFieldTagVersion = "tag_version"       // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt  = "deleted_at"        // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
