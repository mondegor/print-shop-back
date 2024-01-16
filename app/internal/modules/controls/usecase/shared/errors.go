package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrElementTemplateNotFound = NewFactory(
		"errControlsElementTemplateNotFound", ErrorKindUser, "element template with ID={{ .id }} not found")

	FactoryErrFormDataNotFound = NewFactory(
		"errControlsFormDataNotFound", ErrorKindUser, "form with ID={{ .id }} not found")

	FactoryErrFormDataParamNameAlreadyExists = NewFactory(
		"errControlsFormDataParamNameAlreadyExists", ErrorKindUser, "param name '{{ .name }}' already exists")

	FactoryErrFormElementNotFound = NewFactory(
		"errControlsFormElementNotFound", ErrorKindUser, "form element with ID={{ .id }} not found")

	FactoryErrFormElementParamNameAlreadyExists = NewFactory(
		"errControlsFormElementParamNameAlreadyExists", ErrorKindUser, "param name '{{ .name }}' already exists")
)
