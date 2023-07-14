package usecase

import "print-shop-back/pkg/mrerr"

var (
    ErrFormFieldItemTemplateNotFound = mrerr.NewFactory(
        "errFormFieldItemTemplateNotFound", mrerr.ErrorKindUser, "template with ID={{ .id }} not found")

    ErrFormFieldItemParamNameAlreadyExists = mrerr.NewFactory(
        "errFormFieldItemParamNameAlreadyExists", mrerr.ErrorKindUser, "param name '{{ .name }}' is already exists")
)
