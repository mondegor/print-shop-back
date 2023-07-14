package mrerr

type Helper struct {
}

func NewHelper() *Helper {
    return &Helper{}
}

func (h *Helper) WrapErrorForSelect(err error, entityName string) error {
    if ErrStorageNoRowFound.Is(err) {
        return ErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return ErrServiceEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
}

func (h *Helper) WrapErrorForUpdate(err error, entityName string) error {
    if ErrStorageRowsNotAffected.Is(err) {
        return ErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return ErrServiceEntityNotUpdated.Caller(1).Wrap(err, entityName)
}

func (h *Helper) WrapErrorForRemove(err error, entityName string) error {
    if ErrStorageRowsNotAffected.Is(err) {
        return ErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return ErrServiceEntityNotRemoved.Caller(1).Wrap(err, entityName)
}

func (h *Helper) ReturnErrorIfItemNotFound(err error, entityName string) error {
    if err != nil {
        if ErrStorageNoRowFound.Is(err) {
            return ErrServiceEntityNotFound.Wrap(err, entityName)
        }

        return ErrServiceEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
    }

    return nil
}
