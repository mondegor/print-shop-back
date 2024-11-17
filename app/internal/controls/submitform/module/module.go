package module

const (
	Name       = "Controls.SubmitForm"   // Name - название модуля
	Permission = "modControlsSubmitForm" // Permission - разрешение модуля

	DBSchema                      = "printshop_controls"               // DBSchema - схема БД используемая модулем
	DBTableNameElementTemplates   = DBSchema + ".element_templates"    // DBTableNameElementTemplates - таблица БД шаблонов элементов форм
	DBTableNameSubmitForms        = DBSchema + ".submit_forms"         // DBTableNameSubmitForms - таблица БД пользовательских форм
	DBTableNameSubmitFormElements = DBSchema + ".submit_form_elements" // DBTableNameSubmitFormElements - таблица БД элементов форм
	DBTableNameSubmitFormVersions = DBSchema + ".submit_form_versions" // DBTableNameSubmitFormVersions - таблица БД истории форм
	DBFieldTagVersion             = "tag_version"                      // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt              = "deleted_at"                       // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена

	JsonFileNamePattern = "form-%s-%d.json" // JsonFileNamePattern - шаблон названия json файла для выгрузки
	JsonPrettyIndent    = "  "              // JsonPrettyIndent - отступ, который используется при формировании json файлов для скачивания

	UnitSubmitFormPermission  = Permission // UnitSubmitFormPermission - разрешение юнита SubmitForm
	UnitFormElementPermission = Permission // UnitFormElementPermission - разрешение юнита FormElement
)
