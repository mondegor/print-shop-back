package factory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/eventemitter"
	"github.com/mondegor/go-webcore/mrserver/mrprometheus"

	"github.com/mondegor/print-shop-back/internal/app"
)

// NewEventEmitter - создаёт объект eventemitter.Emitter.
func NewEventEmitter(opts app.Options) *eventemitter.Emitter {
	observeEvent := mrprometheus.NewObserveEvent("rest_api", "go")

	opts.Prometheus.MustRegister(
		observeEvent.Collectors()...,
	)

	return eventemitter.New(
		mrsender.EventReceiveFunc(
			func(ctx context.Context, eventName, source string, object any) {
				observeEvent.IncrementEvent(eventName, source)

				mrlog.Ctx(ctx).
					Info().
					Str("event", eventName).
					Str("source", source).
					Any("object", object).
					Send()
			},
		),
	)
}
