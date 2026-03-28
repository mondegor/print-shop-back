package module

const (
	// Name - название модуля.
	Name = "Warehousing"

	// Permission - разрешение модуля.
	Permission = "modWarehousing"

	// GroupContainersMax - максимальное кол-во контейнеров в группе.
	GroupContainersMax = 256

	// DBSchema - схема БД используемая модулем.
	DBSchema = "warehousing"

	// DBTableNameAccountSequences - таблица БД используемая модулем.
	DBTableNameAccountSequences = DBSchema + ".account_sequences"

	// DBTableNameContainers - таблица БД используемая модулем.
	DBTableNameContainers = DBSchema + ".containers"

	// DBTableNameTerritories - таблица БД используемая модулем.
	DBTableNameTerritories = DBSchema + ".territories"

	// DBTableNameStores - таблица БД используемая модулем.
	DBTableNameStores = DBSchema + ".stores"

	// DBTableNameStocks - таблица БД используемая модулем.
	DBTableNameStocks = DBSchema + ".container_stocks"

	// DBFieldTagVersion - поле для хранения версии записи.
	DBFieldTagVersion = "tag_version"

	// DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена.
	DBFieldDeletedAt = "deleted_at"
)
