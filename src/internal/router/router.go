package router

import (
	"avito-test/internal/config"
	"avito-test/internal/service"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Router struct {
	router         *router.Router
	srv            *fasthttp.Server
	logger         *logrus.Logger
	billingService service.Service
	host           string
}

func NewRouter(logger *logrus.Logger, cfg *config.Config,
	billingService service.Service,

) *Router {
	router := router.New()
	srv := &fasthttp.Server{}
	r := &Router{
		router:         router,
		logger:         logger,
		billingService: billingService,
		host:           cfg.Service.Host,
		srv:            srv,
	}
	srv.Handler = r.loggerDecorator(router.Handler)

	registerTenderApi(r)
	registerBidsApi(r)

	r.router.GET("/api/ping", statusHandler)
	return r
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe(r.host)
}

func (r *Router) Shutdown() error {
	return r.srv.Shutdown()
}

func statusHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("ok"))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (r *Router) loggerDecorator(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if recover := recover(); recover != nil {
				r.logger.Println("Recovered in f", recover)
				internalServerErrorResponce(ctx)
			}
		}()
		handler(ctx)
		r.logger.Printf("api request: %s ;status code: %d", ctx.Path(), ctx.Response.StatusCode())
	}
}
