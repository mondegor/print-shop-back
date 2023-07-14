package mrlib

import (
    "fmt"
    "print-shop-back/pkg/mrerr"
)

type (
    UserError struct {
        Id string
        Err *mrerr.AppError
    }

    UserErrorList []UserError
)

func NewUserError(id string, err error) UserError {
    appArr, ok := err.(*mrerr.AppError)

    if !ok {
        appArr = mrerr.New(
            "errMessageForUser",
            err.Error(),
        )
    }

    return UserError{Id: id, Err: appArr}
}

func NewUserErrorList(items ...UserError) *UserErrorList {
    if len(items) > 0 {
        list := append(UserErrorList{}, items...)
        return &list
    }

    return &UserErrorList{}
}

func NewUserErrorListWithError(id string, err error) *UserErrorList {
    return &UserErrorList{NewUserError(id, err)}
}

func (e *UserErrorList) Add(id string, err error) {
    *e = append(*e, NewUserError(id, err))
}

func (e *UserErrorList) Error() string {
    return fmt.Sprintf("%+v", *e)
}
