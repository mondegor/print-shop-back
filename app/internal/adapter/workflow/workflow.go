package workflow

import (
	"github.com/mondegor/go-core/mrworkflow"
	"github.com/mondegor/go-core/mrworkflow/itemstatus"
)

type (
	// FlowMap - интерфейс управления переходами между статусами.
	// Определяет, какие статусы зарегистрированы и какие переходы между ними допустимы.
	FlowMap[Status ~uint8] = mrworkflow.FlowMap[Status]

	// FlowNode - описывает допустимые переходы из одного статуса в другие.
	// From - исходный статус.
	// To - список статусов, в которые разрешён переход из From.
	FlowNode[Status ~uint8] = mrworkflow.FlowNode[Status]

	// ItemStatus - статус элемента (сущности с жизненным циклом).
	// Поддерживаемые статусы: Draft, Enabled, Disabled.
	ItemStatus = itemstatus.Enum
)

// NewFlowMap - создаёт карту допустимых переходов между статусами.
// Параметр list - список узлов переходов (FlowNode), определяющих граф состояний.
// Автоматически строит двунаправленную карту: from→to и to→from.
func NewFlowMap[Status ~uint8](list []FlowNode[Status]) FlowMap[Status] {
	return mrworkflow.NewFlowMap(list)
}
