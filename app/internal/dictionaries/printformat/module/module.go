package module

const (
	Name       = "Dictionaries.PrintFormat"   // Name - название модуля
	Permission = "modDictionariesPrintFormat" // Permission - разрешение модуля

	DBSchema                = "printshop_dictionaries"    // DBSchema - схема БД используемая модулем
	DBTableNamePrintFormats = DBSchema + ".print_formats" // DBTableNamePrintFormats - таблица БД используемая модулем
	DBFieldTagVersion       = "tag_version"               // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt        = "deleted_at"                // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
