package flow

import (
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"
)

// PublicStatusFlow - возвращает стандартную карту возможных переходов PublicStatus.
func PublicStatusFlow() *mrflow.StatusFlow {
	return mrflow.NewStatusFlow(
		[]mrflow.StatusFlowItem{
			{
				From: enum.PublicStatusDraft,
				To: []mrstatus.Getter{
					enum.PublicStatusPublished,
					enum.PublicStatusPublishedShared,
				},
			},
			{
				From: enum.PublicStatusHidden,
				To: []mrstatus.Getter{
					enum.PublicStatusPublished,
					enum.PublicStatusPublishedShared,
				},
			},
			{
				From: enum.PublicStatusPublished,
				To: []mrstatus.Getter{
					enum.PublicStatusHidden,
					enum.PublicStatusPublishedShared,
				},
			},
			{
				From: enum.PublicStatusPublishedShared,
				To: []mrstatus.Getter{
					enum.PublicStatusHidden,
					enum.PublicStatusPublished,
				},
			},
		},
	)
}
