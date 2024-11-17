package module

const (
	Name       = "Controls.ElementTemplate"   // Name - название модуля
	Permission = "modControlsElementTemplate" // Permission - разрешение модуля

	DBSchema                    = "printshop_controls"            // DBSchema - схема БД используемая модулем
	DBTableNameElementTemplates = DBSchema + ".element_templates" // DBTableNameElementTemplates - таблица БД используемая модулем
	DBFieldTagVersion           = "tag_version"                   // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt            = "deleted_at"                    // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена

	JsonFileNamePattern = "controls-template-%d.json" // JsonFileNamePattern - comment const
	JsonPrettyIndent    = "  "                        // JsonPrettyIndent - comment const
)
