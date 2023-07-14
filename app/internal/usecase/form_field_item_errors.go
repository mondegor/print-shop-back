package usecase

import "calc-user-data-back-adm/pkg/mrerr"

var (
    ErrFormFieldItemTemplateNotFound = mrerr.NewFactory(
        "errFormFieldItemTemplateNotFound", mrerr.ErrorKindUser, "template with ID={{ .id }} not found")

    ErrFormFieldItemParamNameAlreadyExists = mrerr.NewFactory(
        "errFormFieldItemParamNameAlreadyExists", mrerr.ErrorKindUser, "param name '{{ .name }}' is already exists")
)
