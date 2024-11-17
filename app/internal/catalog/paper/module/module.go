package module

const (
	Name       = "Catalog.Paper"   // Name - название модуля
	Permission = "modCatalogPaper" // Permission - разрешение модуля

	DBSchema          = "printshop_catalog"  // DBSchema - схема БД используемая модулем
	DBTableNamePapers = DBSchema + ".papers" // DBTableNamePapers - таблица БД используемая модулем
	DBFieldTagVersion = "tag_version"        // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt  = "deleted_at"         // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
