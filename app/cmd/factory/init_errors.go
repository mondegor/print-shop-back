package factory

import (
	"github.com/mondegor/go-sysmess/mrcaller"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/features"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver/mrreq"

	"github.com/mondegor/print-shop-back/internal/app"
)

// InitProtoAppErrors - инициализирует работу с ProtoAppError ошибками.
func InitProtoAppErrors(opts app.Options) {
	var (
		options         []mrerr.ProtoOption
		callerOption    mrerr.ProtoOption
		onCreatedOption mrerr.ProtoOption
	)

	if opts.Cfg.Debugging.ErrorCaller.Enable {
		caller := mrcaller.New(
			mrcaller.WithDepth(opts.Cfg.Debugging.ErrorCaller.Depth),
			mrcaller.WithShowFuncName(opts.Cfg.Debugging.ErrorCaller.ShowFuncName),
			mrcaller.WithFilterStackTrace(
				mrcaller.FilterStackTraceTrimUpper(opts.Cfg.Debugging.ErrorCaller.UpperBounds),
			),
		)

		callerOption = mrerr.WithProtoCaller(
			func() mrerr.StackTracer {
				return caller.StackTrace()
			},
		)

		onCreatedOption = mrerr.WithProtoOnCreated(
			func(_ *mrerr.AppError) (instanceID string) {
				return features.GenerateInstanceID()
			},
		)

		options = append(options, callerOption)
	}

	if opts.Sentry != nil {
		sentry := opts.Sentry
		onCreatedOption = mrerr.WithProtoOnCreated(
			func(err *mrerr.AppError) (instanceID string) {
				if instanceID = sentry.CaptureAppError(err); instanceID != "" {
					return instanceID
				}

				return features.GenerateInstanceID()
			},
		)
	}

	if onCreatedOption != nil {
		options = append(options, onCreatedOption)
	}

	// если ни одна опция не указана, то вызывается инициализация ошибок
	// без использования обработчиков опций по умолчанию
	if len(options) == 0 {
		mrerr.InitDefaultOptions()

		return
	}

	// список ошибок, у которых должны быть индивидуальные опции
	specialErrors := []mrinit.ErrorSettings{
		mrinit.WithOnCreated(mrcore.ErrUnexpectedInternal),
		mrinit.AllDisabled(mrcore.ErrStorageNoRowFound),
		mrinit.AllDisabled(mrcore.ErrStorageRowsNotAffected),
		mrinit.WithCaller(mrcore.ErrUseCaseIncorrectInputData),
		mrinit.AllDisabled(mrcore.ErrHttpRequestParseData),
		mrinit.AllDisabled(mrreq.ErrHttpRequestCorrelationID),
		mrinit.AllDisabled(mrreq.ErrHttpRequestUserIP),
		mrinit.AllDisabled(mrreq.ErrHttpRequestParseUserIP),
	}

	// формируется сопоставление кода ошибки и соответствующих опций по умолчанию
	errorOptionsMap := mrinit.CreateErrorOptionsMap(specialErrors, callerOption, onCreatedOption)

	// обработчик опций по умолчанию, который будет использоваться при вызове mrerr.NewProto()
	defaultHandler := func(_ string, kind mrerr.ErrorKind) []mrerr.ProtoOption {
		if kind == mrerr.ErrorKindUser {
			return nil
		}

		return options
	}

	mrerr.InitDefaultOptions(
		mrerr.ProtoOptionsHandlerFunc(func(code string, kind mrerr.ErrorKind) []mrerr.ProtoOption {
			if op, ok := errorOptionsMap[code]; ok {
				return op
			}

			return defaultHandler(code, kind)
		}),
		mrerr.ProtoOptionsHandlerFunc(defaultHandler),
	)
}
