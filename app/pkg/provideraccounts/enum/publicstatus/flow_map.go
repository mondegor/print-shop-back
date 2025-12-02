package publicstatus

import (
	"github.com/mondegor/go-sysmess/mrstatus"
)

// NewFlowMap - возвращает карту возможных переходов PublicStatus.
func NewFlowMap() mrstatus.FlowMap[Enum] {
	return mrstatus.NewFlowMap(
		[]mrstatus.FlowItem[Enum]{
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
