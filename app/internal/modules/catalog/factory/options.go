package factory

import (
	view_shared "print-shop-back/internal/modules/catalog/controller/http_v1/shared/view"
	dictionaries "print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger          mrcore.Logger
		EventBox        mrcore.EventBox
		ServiceHelper   *mrtool.ServiceHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		RequestParser   *view_shared.Parser
		ResponseSender  *mrresponse.Sender

		LaminateTypeAPI dictionaries.LaminateTypeAPI
		PaperColorAPI   dictionaries.PaperColorAPI
		PaperFactureAPI dictionaries.PaperFactureAPI

		UnitBox      *UnitBoxOptions
		UnitLaminate *UnitLaminateOptions
		UnitPaper    *UnitPaperOptions
	}

	UnitBoxOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}

	UnitLaminateOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}

	UnitPaperOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}
)
