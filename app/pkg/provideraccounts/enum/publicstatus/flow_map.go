package publicstatus

import (
	"print-shop-back/internal/adapter/workflow"
)

// NewFlowMap - возвращает карту возможных переходов PublicStatus.
func NewFlowMap() workflow.FlowMap[Enum] {
	return workflow.NewFlowMap(
		[]workflow.FlowNode[Enum]{
			{
				From: Draft,
				To: []Enum{
					Published,
					PublishedShared,
				},
			},
			{
				From: Hidden,
				To: []Enum{
					Published,
					PublishedShared,
				},
			},
			{
				From: Published,
				To: []Enum{
					Hidden,
					PublishedShared,
				},
			},
			{
				From: PublishedShared,
				To: []Enum{
					Hidden,
					Published,
				},
			},
		},
	)
}
