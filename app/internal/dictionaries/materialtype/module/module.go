package module

const (
	// Name - название модуля.
	Name = "Dictionaries.MaterialType"

	// Permission - разрешение модуля.
	Permission = "modDictionariesMaterialType"

	// LocaleDomain - домен локализации записей.
	LocaleDomain = "dictionaries.material-type"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_dictionaries"

	// DBTableNameMaterialTypes - таблица БД используемая модулем.
	DBTableNameMaterialTypes = DBSchema + ".material_types"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
