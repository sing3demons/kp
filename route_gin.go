package kp

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type ginApplication struct {
	router      *gin.Engine
	middlewares []Middleware
	cfg         Config
}

func newGinServer(cfg Config) IApplication {
	mode := os.Getenv("GIN_MODE")

	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(gin.Recovery())

	return &ginApplication{
		router: app,
		cfg:    cfg,
	}
}

func (app *ginApplication) Consume(topic string, h ServiceHandleFunc) error {
	c := newConsumer(app.cfg)
	return c.Consume(topic, h)
}

func (app *ginApplication) Get(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.GET(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig))
	})
}

func (app *ginApplication) Post(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.POST(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig))
	})
}

func (app *ginApplication) Put(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.PUT(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig))
	})
}

func (app *ginApplication) Delete(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.DELETE(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig))
	})
}

func (app *ginApplication) Patch(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.PATCH(path, func(c *gin.Context) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newGinContext(c, &app.cfg.KafkaConfig))
	})
}

func (app *ginApplication) Use(middlewares ...Middleware) {
	app.middlewares = append(app.middlewares, middlewares...)
}

func (app *ginApplication) Start() {
	srv := http.Server{
		Addr:    ":" + app.cfg.AppConfig.Port,
		Handler: app.router,
	}

	go func() {
		log.Printf("server started on port %s\n", app.cfg.AppConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	log.Println("shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("server exiting")
}
