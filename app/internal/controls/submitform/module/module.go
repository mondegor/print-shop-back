package module

const (
	// Name - название модуля.
	Name = "Controls.SubmitForm"

	// Permission - разрешение модуля.
	Permission = "modControlsSubmitForm"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_controls"

	// DBTableNameElementTemplates - таблица БД шаблонов элементов форм.
	DBTableNameElementTemplates = DBSchema + ".element_templates"

	// DBTableNameSubmitForms - таблица БД пользовательских форм.
	DBTableNameSubmitForms = DBSchema + ".submit_forms"

	// DBTableNameSubmitFormElements - таблица БД элементов форм.
	DBTableNameSubmitFormElements = DBSchema + ".submit_form_elements"

	// DBTableNameSubmitFormVersions - таблица БД истории форм.
	DBTableNameSubmitFormVersions = DBSchema + ".submit_form_versions"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"

	// JsonFileNamePattern - шаблон названия json файла для выгрузки.
	JsonFileNamePattern = "form-%s-%d.json"

	// JsonPrettyIndent - отступ, который используется при формировании json файлов для скачивания.
	JsonPrettyIndent = "  "

	// UnitSubmitFormPermission - разрешение юнита SubmitForm.
	UnitSubmitFormPermission = Permission

	// UnitFormElementPermission - разрешение юнита FormElement.
	UnitFormElementPermission = Permission
)
