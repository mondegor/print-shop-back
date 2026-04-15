package factory

import (
	"context"

	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver/mrprometheus"

	"github.com/mondegor/print-shop-back/internal/app"
)

// InitEventEmitter - создаёт объект mrevent.Emitter.
func InitEventEmitter(opts app.Options) mrevent.Emitter {
	observeEvent := mrprometheus.NewObserveEvent("rest_api", "go")

	opts.Prometheus.Add(observeEvent.Collectors()...)

	return mrevent.NewEmitter(
		mrevent.ReceiveFunc(
			func(ctx context.Context, eventName string, args ...any) {
				source, eventName := mrevent.ExtractEventName(eventName)

				observeEvent.IncrementEvent(eventName, source)

				opts.Logger.Info(
					ctx, "EventEmitter",
					append(
						args,
						"event", eventName,
						"source", source,
					)...,
				)
			},
		),
	)
}
