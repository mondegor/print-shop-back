package module

const (
	// Name - название модуля.
	Name = "Dictionaries.PaperFacture"

	// Permission - разрешение модуля.
	Permission = "modDictionariesPaperFacture"

	// LocaleDomain - домен локализации записей.
	LocaleDomain = "dictionaries.paper-facture"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_dictionaries"

	// DBTableNamePaperFactures - таблица БД используемая модулем.
	DBTableNamePaperFactures = DBSchema + ".paper_factures"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
