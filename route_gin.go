package kp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginApplication struct {
	router      *gin.Engine
	middlewares []Middleware
	cfg         *Config
	log         ILogger
}

func newGinServer(cfg *Config, log ILogger) IRouter {
	app := gin.New()
	app.Use(gin.Recovery())

	return &ginApplication{
		router: app,
		cfg:    cfg,
		log:    log,
	}
}

func (app *ginApplication) Get(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.GET(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *ginApplication) Post(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.POST(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *ginApplication) Put(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.PUT(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *ginApplication) Delete(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.DELETE(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *ginApplication) Patch(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.PATCH(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *ginApplication) Use(middlewares ...Middleware) {
	app.middlewares = append(app.middlewares, middlewares...)
}

func (app *ginApplication) Register() *http.Server {
	srv := &http.Server{
		Addr:    ":" + app.cfg.AppConfig.Port,
		Handler: app.router,
	}

	return srv
	// go func() {
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatal(err)
	// 	}
	// }()

	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// <-shutdown
	// log.Println("shutting down...")
	// if err := srv.Shutdown(context.Background()); err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("server exiting")
}
