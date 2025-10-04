package module

const (
	// Name - comment const.
	Name = "Calculations.QueryHistory"

	// Permission - разрешение модуля.
	Permission = "modCalculationsQueryHistory"

	// DBSchema - схема БД используемая модулем.
	DBSchema = "printshop_calculation"

	// DBTableNameQueryHistory - таблица БД используемая модулем.
	DBTableNameQueryHistory = DBSchema + ".query_history"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
