package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/errors"

	"print-shop-back/internal/warehousing/enum/locationkind"
)

const (
	containerPrefix = "C/"
	groupPrefix     = "G/"
)

type (
	// AccountSequence - comment struct.
	AccountSequence struct {
		storageSequence accountSequenceStorage
		errorWrapper    errors.Wrapper
	}

	accountSequenceStorage interface {
		NextContainerID(ctx context.Context, accountID uuid.UUID) (id uint64, err error)
	}
)

// NewAccountSequence - создаёт объект AccountSequence.
func NewAccountSequence(
	storageSequence accountSequenceStorage,
) *AccountSequence {
	return &AccountSequence{
		storageSequence: storageSequence,
		errorWrapper:    errors.NewServiceOperationFailedWrapper(),
	}
}

// ContainerCode - comment method.
func (uc *AccountSequence) ContainerCode(ctx context.Context, accountID uuid.UUID, kind locationkind.Enum) (string, error) {
	var prefix string

	switch kind {
	case locationkind.Group:
		prefix = groupPrefix
	case locationkind.Container:
		prefix = containerPrefix
	default:
		return "", errors.ErrInternalIncorrectInputData.WithDetails("unexpected location kind", "kind", kind)
	}

	id, err := uc.storageSequence.NextContainerID(ctx, accountID)
	if err != nil {
		return "", uc.errorWrapper.Wrap(err)
	}

	return prefix + strings.ToUpper(strconv.FormatUint(id, 36)), nil
}
