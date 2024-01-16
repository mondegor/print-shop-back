package factory

import (
	catalog "print-shop-back/pkg/modules/dictionaries"

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

		LaminateTypeAPI catalog.LaminateTypeAPI
		PaperColorAPI   catalog.PaperColorAPI
		PaperFactureAPI catalog.PaperFactureAPI

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
