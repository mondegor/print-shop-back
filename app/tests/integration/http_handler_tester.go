package integration

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-storage/mrtests/infra"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtests/helpers"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/print-shop-back/cmd/factory"
	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/tests"
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

	ctx := helpers.ContextWithNopLogger()
	cfg, err := config.Create(
		config.Args{
			WorkDir:    tests.AppWorkDir(),
			DotEnvPath: tests.AppDotEnvPathForTests(),
			Stdout:     os.Stdout,
		},
	)
	require.NoError(t, err)

	pgt := infra.NewPostgresTester(t, tests.DBSchemas(), tests.ExcludedDBTables())
	pgt.ApplyMigrations(tests.AppMigrationsDir())

	rds := infra.NewRedisTester(t)

	fpool, err := factory.NewFileProviderPool(ctx, cfg)
	require.NoError(t, err)

	ctx, opts, err := factory.InitAppEnvironment(
		ctx,
		app.Options{
			Cfg:                 cfg,
			PostgresConnManager: pgt.ConnManager(),
			RedisAdapter:        rds.Conn(),
			FileProviderPool:    fpool,
		},
	)
	require.NoError(t, err)

	router, err := factory.NewRestRouterWithRegisterHandlers(ctx, opts)
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
	mrlib.CallEachFunc(t.ctx, t.opts.OpenedResources)
	t.pgt.Destroy(t.ctx)
	t.rds.Destroy(t.ctx)

	err := t.fpool.Close()
	require.NoError(t.parentT, err)
}
