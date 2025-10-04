package module

const (
	// Name - название модуля.
	Name = "Dictionaries.PrintFormat"

	// Permission - разрешение модуля.
	Permission = "modDictionariesPrintFormat"

	// LocaleDomain - домен локализации записей.
	LocaleDomain = "dictionaries.print-format"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_dictionaries"

	// DBTableNamePrintFormats - таблица БД используемая модулем.
	DBTableNamePrintFormats = DBSchema + ".print_formats"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
