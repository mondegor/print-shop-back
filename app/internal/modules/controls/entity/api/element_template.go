package entity_api

import (
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplateHead struct { // DB: printdata_controls.element_templates
		ID        mrtype.KeyInt32
		ParamName string
		Caption   string
	}
)
