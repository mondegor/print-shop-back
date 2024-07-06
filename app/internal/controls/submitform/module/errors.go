package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

var (
	// ErrSubmitFormRequired - form ID is required.
	ErrSubmitFormRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsSubmitFormRequired", mrerr.ErrorKindUser, "form ID is required")

	// ErrSubmitFormNotFound - form with ID not found.
	ErrSubmitFormNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsSubmitFormNotFound", mrerr.ErrorKindUser, "form with ID={{ .id }} not found")

	// ErrSubmitFormRewriteNameAlreadyExists - rewrite name already exists.
	ErrSubmitFormRewriteNameAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsSubmitFormRewriteNameAlreadyExists", mrerr.ErrorKindUser, "rewrite name '{{ .name }}' already exists")

	// ErrSubmitFormParamNameAlreadyExists - param name already exists.
	ErrSubmitFormParamNameAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsSubmitFormParamNameAlreadyExists", mrerr.ErrorKindUser, "param name '{{ .name }}' already exists")

	// ErrSubmitFormIsDisabled - form with ID is disabled.
	ErrSubmitFormIsDisabled = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsSubmitFormIsDisabled", mrerr.ErrorKindUser, "form with ID={{ .id }} is disabled")

	// ErrFormElementNotFound - form element with ID not found.
	ErrFormElementNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsFormElementNotFound", mrerr.ErrorKindUser, "form element with ID={{ .id }} not found")

	// ErrFormElementParamNameAlreadyExists - param name already exists.
	ErrFormElementParamNameAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsFormElementParamNameAlreadyExists", mrerr.ErrorKindUser, "param name '{{ .name }}' already exists")

	// ErrFormElementDetailingNotAllowed - item detailing not allowed for form detailing.
	ErrFormElementDetailingNotAllowed = mrerrfactory.NewProtoAppErrorByDefault(
		"errControlsFormElementDetailingNotAllowed", mrerr.ErrorKindUser, "item detailing '{{ .name1 }}' not allowed for form detailing '{{ .name2 }}'")
)

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrSubmitFormRequired,
		ErrSubmitFormNotFound,
		ErrSubmitFormRewriteNameAlreadyExists,
		ErrSubmitFormParamNameAlreadyExists,
		ErrSubmitFormIsDisabled,
		ErrFormElementNotFound,
		ErrFormElementParamNameAlreadyExists,
		ErrFormElementDetailingNotAllowed,
	}
}
