package module

const (
	Name       = "Dictionaries.PaperColor"   // Name - название модуля
	Permission = "modDictionariesPaperColor" // Permission - разрешение модуля

	DBSchema               = "printshop_dictionaries"   // DBSchema - схема БД используемая модулем
	DBTableNamePaperColors = DBSchema + ".paper_colors" // DBTableNamePaperColors - таблица БД используемая модулем
	DBFieldTagVersion      = "tag_version"              // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt       = "deleted_at"               // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
