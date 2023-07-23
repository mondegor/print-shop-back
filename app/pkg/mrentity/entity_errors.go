package mrentity

import "print-shop-back/pkg/mrerr"

var (
    ErrInternalListOfFieldsIsEmpty = mrerr.NewFactory(
        "errInternalListOfFieldsIsEmpty", mrerr.ErrorKindInternalNotice, "the list of fields is empty")
)
