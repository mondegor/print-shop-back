package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrSubmitFormRequired = NewFactory(
		"errControlsSubmitFormRequired", ErrorKindUser, "form ID is required")

	FactoryErrSubmitFormNotFound = NewFactory(
		"errControlsSubmitFormNotFound", ErrorKindUser, "form with ID={{ .id }} not found")

	FactoryErrSubmitFormRewriteNameAlreadyExists = NewFactory(
		"errControlsSubmitFormRewriteNameAlreadyExists", ErrorKindUser, "rewrite name '{{ .name }}' already exists")

	FactoryErrSubmitFormParamNameAlreadyExists = NewFactory(
		"errControlsSubmitFormParamNameAlreadyExists", ErrorKindUser, "param name '{{ .name }}' already exists")

	FactoryErrSubmitFormIsDisabled = NewFactory(
		"errControlsSubmitFormIsDisabled", ErrorKindUser, "form with ID={{ .id }} is disabled")

	FactoryErrFormElementNotFound = NewFactory(
		"errControlsFormElementNotFound", ErrorKindUser, "form element with ID={{ .id }} not found")

	FactoryErrFormElementParamNameAlreadyExists = NewFactory(
		"errControlsFormElementParamNameAlreadyExists", ErrorKindUser, "param name '{{ .name }}' already exists")

	FactoryErrFormElementDetailingNotAllowed = NewFactory(
		"errControlsFormElementDetailingNotAllowed", ErrorKindUser, "item detailing '{{ .name1 }}' not allowed for form detailing '{{ .name2 }}'")
)
