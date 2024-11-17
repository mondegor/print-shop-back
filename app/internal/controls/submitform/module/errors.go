package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrSubmitFormRequired - form ID is required.
	ErrSubmitFormRequired = mrerr.NewProto(
		"controls.errSubmitFormRequired", mrerr.ErrorKindUser, "form ID is required")

	// ErrSubmitFormNotFound - form with ID not found.
	ErrSubmitFormNotFound = mrerr.NewProto(
		"controls.errSubmitFormNotFound", mrerr.ErrorKindUser, "form with ID={{ .id }} not found")

	// ErrSubmitFormRewriteNameAlreadyExists - rewrite name already exists.
	ErrSubmitFormRewriteNameAlreadyExists = mrerr.NewProto(
		"controls.errSubmitFormRewriteNameAlreadyExists", mrerr.ErrorKindUser, "rewrite name '{{ .name }}' already exists")

	// ErrSubmitFormParamNameAlreadyExists - param name already exists.
	ErrSubmitFormParamNameAlreadyExists = mrerr.NewProto(
		"controls.errSubmitFormParamNameAlreadyExists", mrerr.ErrorKindUser, "param name '{{ .name }}' already exists")

	// ErrSubmitFormIsDisabled - form with ID is disabled.
	ErrSubmitFormIsDisabled = mrerr.NewProto(
		"controls.errSubmitFormIsDisabled", mrerr.ErrorKindUser, "form with ID={{ .id }} is disabled")

	// ErrFormElementNotFound - form element with ID not found.
	ErrFormElementNotFound = mrerr.NewProto(
		"controls.errFormElementNotFound", mrerr.ErrorKindUser, "form element with ID={{ .id }} not found")

	// ErrFormElementParamNameAlreadyExists - param name already exists.
	ErrFormElementParamNameAlreadyExists = mrerr.NewProto(
		"controls.errFormElementParamNameAlreadyExists", mrerr.ErrorKindUser, "param name '{{ .name }}' already exists")

	// ErrFormElementDetailingNotAllowed - item detailing not allowed for form detailing.
	ErrFormElementDetailingNotAllowed = mrerr.NewProto(
		"controls.errFormElementDetailingNotAllowed", mrerr.ErrorKindUser, "item detailing '{{ .name1 }}' not allowed for form detailing '{{ .name2 }}'")
)
