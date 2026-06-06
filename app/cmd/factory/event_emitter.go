package factory

import (
	"context"

	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver/mrprometheus"

	"print-shop-back/internal/app"
)

// InitEventEmitter - создаёт объект mrevent.Emitter.
func InitEventEmitter(opts app.Options) mrevent.Emitter {
	receiveFunc := func(ctx context.Context, eventName string, args ...any) {
		source, eventName := mrevent.ExtractEventName(eventName)

		opts.Logger.Info(
			ctx, "EventEmitter",
			append(
				args,
				"event", eventName,
				"source", source,
			)...,
		)
	}

	if opts.Prometheus != nil {
		observeEvent := mrprometheus.NewObserveEvent("rest_api", "go")
		opts.Prometheus.Add(observeEvent.Collectors()...)

		receiveFunc = func(ctx context.Context, eventName string, args ...any) {
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
		}
	}

	return mrevent.NewEmitter(mrevent.ReceiveFunc(receiveFunc))
}
