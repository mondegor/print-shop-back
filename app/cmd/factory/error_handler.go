package factory

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/config"
)

// NewErrorHandler - создаёт объект mrapp.AppErrorHandler.
func NewErrorHandler(_ mrlog.Logger, _ config.Config) *mrapp.ErrorHandler {
	return mrapp.NewErrorHandlerWithHook(
		func(ctx context.Context, errType mrcore.AnalyzedErrorType, err *mrerr.AppError) {
			if errType == mrcore.AnalyzedErrorTypeUser || errType == mrcore.AnalyzedErrorTypeProtoUser {
				// 1. пользовательские ошибки: AppError + User, ProtoAppError + User
				mrlog.Ctx(ctx).Debug().Err(err).Send()
			} else {
				// 2. системные ошибки: AppError + Internal/System;
				// 3. ProtoAppError + Internal/System (нужно найти место их создания и добавить у для них вызов New()/Wrap());
				// 4. остальные ошибки: которые не были обёрнуты в ProtoAppError (нужно найти место их создания и обернуть);
				mrlog.Ctx(ctx).Error().Err(err).Send()
			}
		},
	)
}
