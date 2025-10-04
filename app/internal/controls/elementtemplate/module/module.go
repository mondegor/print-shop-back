package module

const (
	// Name - название модуля.
	Name = "Controls.ElementTemplate"

	// Permission - разрешение модуля.
	Permission = "modControlsElementTemplate"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_controls"

	// DBTableNameElementTemplates - таблица БД используемая модулем.
	DBTableNameElementTemplates = DBSchema + ".element_templates"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"

	// JsonFileNamePattern - comment const.
	JsonFileNamePattern = "controls-template-%d.json"

	// JsonPrettyIndent - comment const.
	JsonPrettyIndent = "  "
)
