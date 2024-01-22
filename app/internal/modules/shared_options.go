package modules

import (
	"print-shop-back/config"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Cfg              *config.Config
		Logger           mrcore.Logger
		EventBox         mrcore.EventBox
		ServiceHelper    *mrtool.ServiceHelper
		PostgresAdapter  *mrpostgres.ConnAdapter
		RedisAdapter     *mrredis.ConnAdapter
		FileProviderPool *mrstorage.FileProviderPool
		Locker           mrcore.Locker
		Translator       *mrlang.Translator
		RequestParsers   *RequestParsers
		ResponseSender   *mrresponse.Sender
		AccessControl    mrcore.AccessControl

		DictionariesLaminateTypeAPI dictionaries.LaminateTypeAPI
		DictionariesPaperColorAPI   dictionaries.PaperColorAPI
		DictionariesPaperFactureAPI dictionaries.PaperFactureAPI
		DictionariesPrintFormatAPI  dictionaries.PrintFormatAPI
		OrdererAPI                  mrorderer.API
	}

	RequestParsers struct {
		Path       mrserver.RequestParserPath
		Base       *mrparser.Base
		ItemStatus *mrparser.ItemStatus
		KeyInt32   *mrparser.KeyInt32
		SortPage   *mrparser.SortPage
		Validator  *mrparser.Validator
	}
)
