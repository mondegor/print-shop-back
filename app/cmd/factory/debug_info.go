package factory

import (
	authdebug "github.com/mondegor/go-components/mrauth/bag/debug"
	"github.com/mondegor/go-components/mrauth/model/secureoperation"
	"github.com/mondegor/go-sysmess/errors/runtime/hint"
)

// InitDebugInfo - возвращает функцию, которая подготавливает
// отладочные данные в случае, если дебаггер режим включён.
func InitDebugInfo(isDebug bool) func(value any) string {
	if !isDebug {
		return func(_ any) string {
			return ""
		}
	}

	return func(value any) string {
		switch o := value.(type) {
		case error:
			return hint.DetailedError(o)
		case secureoperation.SecureOperation:
			return authdebug.Info(o, isDebug)
		default:
			return ""
		}
	}
}
