package mrhttp

import (
    "net/http"
    "print-shop-back/pkg/mrapp"
    "reflect"
    "runtime"

    "github.com/julienschmidt/httprouter"
)

// go get -u github.com/julienschmidt/httprouter

// Make sure the Router conforms with the mrapp.Router interface
var _ mrapp.Router = (*Router)(nil)

type Router struct {
    router *httprouter.Router
    generalHandler http.Handler
    logger mrapp.Logger
    validator mrapp.Validator
}

func NewRouter(logger mrapp.Logger, validator mrapp.Validator) *Router {
    router := httprouter.New()

    // r.GlobalOPTIONS
    // rt.router.MethodNotAllowed
    // rt.router.NotFound

    return &Router{
        router: router,
        generalHandler: router,
        logger: logger,
        validator: validator,
    }
}

func (rt *Router) RegisterMiddleware(handlers ...mrapp.HttpMiddleware) {
    // recursion call: handler1(handler2(handler3(router())))
    for i := len(handlers) - 1; i >= 0; i-- {
        rt.generalHandler = handlers[i].Middleware(rt.generalHandler)
        rt.logger.Info(
            "Registered Middleware %s",
            runtime.FuncForPC(reflect.ValueOf(rt.generalHandler).Pointer()).Name(),
        )
    }
}

func (rt *Router) Register(controllers ...mrapp.HttpController) {
    for _, controller := range controllers {
        controller.AddHandlers(rt)
    }
}

func (rt *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
    rt.router.Handler(method, path, handler)
}

func (rt *Router) HttpHandlerFunc(method, path string, handler mrapp.HttpHandlerFunc) {
    rt.router.Handler(method, path, rt.MiddlewareLast(handler))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rt.generalHandler.ServeHTTP(w, r)
}
