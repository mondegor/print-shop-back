package module

const (
	// Name - название модуля.
	Name = "Dictionaries.PaperColor"

	// Permission - разрешение модуля.
	Permission = "modDictionariesPaperColor"

	// LocaleDomain - домен локализации записей.
	LocaleDomain = "dictionaries.paper-colors"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_dictionaries"

	// DBTableNamePaperColors - таблица БД используемая модулем.
	DBTableNamePaperColors = DBSchema + ".paper_colors"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
