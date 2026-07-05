package trace

import (
	"context"

	"github.com/mondegor/go-sysmess/mrtrace"
	tracectx "github.com/mondegor/go-sysmess/mrtrace/context"
)

const (
	// KeyCorrelationID - название ключа ID корреляции.
	KeyCorrelationID = mrtrace.KeyCorrelationID
)

type (
	// Tracer - интерфейс для трассировки запросов между сервисами.
	// Фиксирует входящие/исходящие запросы и их параметры для аудита и отладки.
	Tracer = mrtrace.Tracer

	// ContextManager - управляет идентификаторами процессов в контексте для трассировки.
	// Позволяет получать, устанавливать и генерировать ID процессов (request_id, process_id и др.).
	ContextManager = mrtrace.ContextManager
)

// WithCorrelationID - возвращает новый контекст с установленным ID корреляции запроса.
// ID корреляции используется для связывания цепочки запросов между сервисами.
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return tracectx.WithCorrelationID(ctx, correlationID)
}

// NopTracer - создаёт Tracer, который игнорирует все данные трассировки.
func NopTracer() Tracer {
	return mrtrace.NopTracer()
}
