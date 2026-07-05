package validate

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver/request/parser"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/pkg/provideraccounts/enum/publicstatus"
)

type (
	// RequestPublicStatusParser - comment interface.
	RequestPublicStatusParser interface { // TODO: ПЕРЕНЕСТИ
		FilterPublicStatusList(r *http.Request, key string) []publicstatus.Enum
	}

	// PublicStatusParser - парсер publicstatus.Enum.
	PublicStatusParser struct {
		*parser.EnumList[publicstatus.Enum]
	}
)

// NewPublicStatusParser - создаёт объект PublicStatusParser.
func NewPublicStatusParser(logger log.Logger) *PublicStatusParser {
	return &PublicStatusParser{
		EnumList: parser.NewEnumList(
			logger,
			publicstatus.ParseList,
		),
	}
}

// NewPublicStatusParserWithDefault - создаёт объект PublicStatusParser со статусами по умолчанию.
func NewPublicStatusParserWithDefault(logger log.Logger, items []publicstatus.Enum) *PublicStatusParser {
	return &PublicStatusParser{
		EnumList: parser.NewEnumListWithDefault(
			logger,
			items,
			publicstatus.ParseList,
		),
	}
}

// FilterPublicStatusList - возвращает массив publicstatus.Enum поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *PublicStatusParser) FilterPublicStatusList(r *http.Request, key string) []publicstatus.Enum {
	return p.FilterEnumList(r, key)
}
