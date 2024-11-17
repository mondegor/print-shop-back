package module

const (
	Name       = "Dictionaries.MaterialType"   // Name - название модуля
	Permission = "modDictionariesMaterialType" // Permission - разрешение модуля

	DBSchema                 = "printshop_dictionaries"     // DBSchema - схема БД используемая модулем
	DBTableNameMaterialTypes = DBSchema + ".material_types" // DBTableNameMaterialTypes - таблица БД используемая модулем
	DBFieldTagVersion        = "tag_version"                // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt         = "deleted_at"                 // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
