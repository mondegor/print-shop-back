package module

const (
	Name       = "Calculations.QueryHistory"   // Name - comment const
	Permission = "modCalculationsQueryHistory" // Permission - разрешение модуля

	DBSchema                = "printshop_calculation"     // DBSchema - схема БД используемая модулем
	DBTableNameQueryHistory = DBSchema + ".query_history" // DBTableNameQueryHistory - таблица БД используемая модулем
	DBFieldDeletedAt        = "deleted_at"                // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
