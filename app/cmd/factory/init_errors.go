package factory

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/instance"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrerr/stacktrace"
	"github.com/mondegor/go-sysmess/mrmsg"

	"github.com/mondegor/print-shop-back/internal/app"
)

// InitProtoAppErrors - инициализирует работу с ProtoAppError ошибками.
func InitProtoAppErrors(opts app.Options) {
	var (
		userOptions     []mrerr.Option
		intSysOptions   []mrerr.Option
		callerOption    mrerr.Option
		onCreatedOption mrerr.Option
	)

	if opts.Cfg.Debugging.ErrorCaller.IsEnabled {
		caller := stacktrace.New(
			stacktrace.WithDepth(opts.Cfg.Debugging.ErrorCaller.Depth),
			stacktrace.WithStackTraceFilter(
				stacktrace.TrimUpperFilter(opts.Cfg.Debugging.ErrorCaller.UpperBounds),
			),
		)

		callerOption = mrerr.WithCaller(
			func() mrerr.StackTracer {
				return caller.StackTrace()
			},
		)

		onCreatedOption = mrerr.WithOnCreated(
			func(context.Context, error) (instanceID string) {
				return instance.GenerateID()
			},
		)

		intSysOptions = append(intSysOptions, callerOption)
	}

	if opts.Sentry != nil {
		sentry := opts.Sentry
		onCreatedOption = mrerr.WithOnCreated(
			func(ctx context.Context, err error) (instanceID string) {
				// ???????????????????????????? ПРОВЕРИТЬ!!!
				if instanceID = sentry.CaptureAppError(ctx, mr.CastOrWrapUnexpectedInternal(err)); instanceID != "" {
					return instanceID
				}

				return instance.GenerateID()
			},
		)
	}

	if onCreatedOption != nil {
		intSysOptions = append(intSysOptions, onCreatedOption)
	}

	messageArgsReplacer := mrerr.WithArgsReplacer(
		func(message string) mrerr.MessageReplacer {
			return mrmsg.NewMessageReplacer("{", "}", message)
		},
	)

	userOptions = append(userOptions, messageArgsReplacer)
	intSysOptions = append(intSysOptions, messageArgsReplacer)

	mrerr.InitDefaultOptions(
		mrerr.OptionsHandlerFunc(func(kind mrerr.ErrorKind, _, _ string) []mrerr.Option {
			// // следующий фрагмент кода генерит вспомогательный код, который необходимо скопировать
			// // в localization/dict/errcat/generate.go для локализации пользовательских ошибок в проекте
			// if kind == mrerr.ErrorKindUser || code != "" {
			// 	mrlog.DebugFunc(
			//	    opts.Logger,
			// 		func() string {
			// 			formattedMessage, placeholders := gotext.NewMessageFormatter("{", "}").Format(message)
			// 			formattedMessage = strconv.Quote(formattedMessage)
			//
			// 			if len(placeholders) > 0 {
			// 				placeholders = mrmsg.NewPlaceholderExtractor("{", "}").Extract(message)
			// 				for i := range placeholders {
			// 					placeholders[i] = strconv.Quote(strings.TrimLeft(strings.TrimRight(placeholders[i], "}"), "{"))
			// 				}
			// 				formattedMessage += ", " + strings.Join(placeholders, ", ")
			// 			}
			//
			// 			if code == "" {
			// 				code = "EMPTY"
			// 			}
			//
			// 			return kind.String() + "|" + code + `: p.Sprintf(` + formattedMessage + `)`
			// 		},
			// 	)
			// }
			if kind == mrerr.ErrorKindUser {
				return userOptions
			}

			return intSysOptions
		}),
	)
}
