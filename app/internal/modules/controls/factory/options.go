package factory

import (
	usecase_api "print-shop-back/internal/modules/controls/usecase/api"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger          mrcore.Logger
		EventBox        mrcore.EventBox
		ServiceHelper   *mrtool.ServiceHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		OrdererAPI      mrorderer.API

		ElementTemplateAPI usecase_api.ElementTemplateAPI

		UnitElementTemplate *UnitElementTemplateOptions
	}

	UnitElementTemplateOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}
)
