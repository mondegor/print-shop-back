package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrElementTemplateNotFound = NewFactory(
		"errControlsElementTemplateNotFound", ErrorKindUser, "element template with ID={{ .id }} not found")

	FactoryErrSubmitFormNotFound = NewFactory(
		"errControlsSubmitFormNotFound", ErrorKindUser, "form with ID={{ .id }} not found")

	FactoryErrSubmitFormParamNameAlreadyExists = NewFactory(
		"errControlsSubmitFormParamNameAlreadyExists", ErrorKindUser, "param name '{{ .name }}' already exists")

	FactoryErrFormElementNotFound = NewFactory(
		"errControlsFormElementNotFound", ErrorKindUser, "form element with ID={{ .id }} not found")

	FactoryErrFormElementParamNameAlreadyExists = NewFactory(
		"errControlsFormElementParamNameAlreadyExists", ErrorKindUser, "param name '{{ .name }}' already exists")
)
