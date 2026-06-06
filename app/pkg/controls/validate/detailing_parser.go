package validate

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver/request/parser"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/pkg/controls/enum/elementdetailing"
)

type (
	// RequestDetailingParser - comment interface.
	RequestDetailingParser interface { // TODO: ПЕРЕНЕСТИ
		FilterElementDetailingList(r *http.Request, key string) []elementdetailing.Enum
	}

	// DetailingParser - парсер elementdetailing.Enum.
	DetailingParser struct {
		*parser.EnumList[elementdetailing.Enum]
	}
)

// NewDetailingParser - создаёт объект DetailingParser.
func NewDetailingParser(logger log.Logger) *DetailingParser {
	return &DetailingParser{
		EnumList: parser.NewEnumList(
			logger,
			elementdetailing.ParseList,
		),
	}
}

// NewDetailingParserWithDefault - создаёт объект DetailingParser со статусами по умолчанию.
func NewDetailingParserWithDefault(logger log.Logger, items []elementdetailing.Enum) *DetailingParser {
	return &DetailingParser{
		EnumList: parser.NewEnumListWithDefault(
			logger,
			items,
			elementdetailing.ParseList,
		),
	}
}

// FilterElementDetailingList - возвращает массив elementdetailing.Enum поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *DetailingParser) FilterElementDetailingList(r *http.Request, key string) []elementdetailing.Enum {
	return p.FilterEnumList(r, key)
}
