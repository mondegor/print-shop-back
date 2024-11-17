package module

const (
	Name       = "Catalog.Laminate"   // Name - название модуля
	Permission = "modCatalogLaminate" // Permission - разрешение модуля

	DBSchema             = "printshop_catalog"     // DBSchema - схема БД используемая модулем
	DBTableNameLaminates = DBSchema + ".laminates" // DBTableNameLaminates - таблица БД используемая модулем
	DBFieldTagVersion    = "tag_version"           // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt     = "deleted_at"            // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
