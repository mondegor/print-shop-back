package module

const (
	Name       = "Dictionaries.PaperFacture"   // Name - название модуля
	Permission = "modDictionariesPaperFacture" // Permission - разрешение модуля

	DBSchema                 = "printshop_dictionaries"     // DBSchema - схема БД используемая модулем
	DBTableNamePaperFactures = DBSchema + ".paper_factures" // DBTableNamePaperFactures - таблица БД используемая модулем
	DBFieldTagVersion        = "tag_version"                // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt         = "deleted_at"                 // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
