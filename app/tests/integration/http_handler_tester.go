package integration

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/mondegor/go-storage/mrtests/infra"
	"github.com/mondegor/go-sysmess/mrstorage"
	"github.com/mondegor/go-sysmess/util/xio"
	"github.com/mondegor/go-sysmess/wire/mrlog"
	"github.com/mondegor/go-sysmess/wire/mrtrace"
	"github.com/mondegor/go-webcore/mrtests/helpers"
	"github.com/stretchr/testify/require"

	"print-shop-back/cmd/factory"
	"print-shop-back/cmd/factory/service/rest"
	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/app"
	"print-shop-back/tests"
)

type (
	// HttpHandlerTester - вспомогательный объект для тестирования http обработчиков.
	HttpHandlerTester struct {
		parentT *testing.T
		ctx     context.Context
		opts    app.Options
		router  http.Handler
		pgt     *infra.PostgresTester
		rds     *infra.RedisTester
		fpool   *mrstorage.FileProviderPool
	}
)

// NewHandlerTester - создаёт объект HttpHandlerTester.
func NewHandlerTester(t *testing.T) *HttpHandlerTester {
	t.Helper()

	ctx := context.Background()
	cfg, err := config.Create(
		config.CmdArgs{
			WorkDir:     tests.AppWorkDir(),
			Environment: "tests",
		},
		os.Stdout,
	)
	require.NoError(t, err)

	logger := log.NopLogger()
	tracer := trace.NopTracer()
	traceManager, err := mrtrace.InitTraceContextManager(mrlog.DefaultProcessIDs(), logger)
	require.NoError(t, err)

	pgt := infra.NewPostgresTester(t, tests.DBSchemas(), tests.ExcludedDBTables())
	pgt.ApplyMigrations(tests.AppMigrationsDir())

	rds := infra.NewRedisTester(t)

	fpool, err := factory.InitFileProviderPool(logger, tracer, cfg)
	require.NoError(t, err)

	opts, err := factory.InitAppEnvironment(
		app.Options{
			Cfg:                 cfg,
			Logger:              logger,
			Tracer:              tracer,
			TraceManager:        traceManager,
			OpenedResources:     xio.NewCloseManager(logger),
			PostgresConnManager: pgt.ConnManager(),
			RedisAdapter:        rds.Conn(),
			FileProviderPool:    fpool,
		},
	)
	require.NoError(t, err)

	router, err := rest.InitRestRouterWithHandlers(opts)
	require.NoError(t, err)

	return &HttpHandlerTester{
		parentT: t,
		ctx:     ctx,
		opts:    opts,
		router:  router,
		pgt:     pgt,
		rds:     rds,
		fpool:   fpool,
	}
}

// Context - возвращает текущий контекст.
func (t *HttpHandlerTester) Context() context.Context {
	return t.ctx
}

// Options - возвращает опции приложения.
func (t *HttpHandlerTester) Options() app.Options {
	return t.opts
}

// Router - возвращает текущий роутер.
func (t *HttpHandlerTester) Router() http.Handler {
	return t.router
}

// ExecRequest - исполняет текущий запрос.
func (t *HttpHandlerTester) ExecRequest(r *helpers.HttpRequest, structResponse any) (statusCode int, err error) {
	return r.Exec(t.router, structResponse)
}

// Clean - очищает ресурсы после завершения тестирования обработчика.
func (t *HttpHandlerTester) Clean() {
	t.opts.OpenedResources.Close()
	t.pgt.Destroy(t.ctx)
	t.rds.Destroy(t.ctx)

	err := t.fpool.Close()
	require.NoError(t.parentT, err)
}
