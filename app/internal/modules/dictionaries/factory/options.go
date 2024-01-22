package factory

import (
	view_shared "print-shop-back/internal/modules/dictionaries/controller/http_v1/shared/view"

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

		UnitLaminateType       *UnitLaminateTypeOptions
		UnitPaperColor         *UnitPaperColorOptions
		UnitPaperFacture       *UnitPaperFactureOptions
		UnitPrintFormatFacture *UnitPrintFormatOptions
	}

	UnitLaminateTypeOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}

	UnitPaperColorOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}

	UnitPaperFactureOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}

	UnitPrintFormatOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}
)
