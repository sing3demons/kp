package kp

import (
	"net/http"
	"time"
)

type httpApplication struct {
	mux         *http.ServeMux
	middlewares []Middleware
	cfg         *Config
	log         ILogger
}

func newServer(cfg *Config, log ILogger) IRouter {
	app := http.NewServeMux()

	return &httpApplication{
		mux: app,
		cfg: cfg,
		log: log,
	}
}

func (app *httpApplication) Use(middlewares ...Middleware) {
	app.middlewares = append(app.middlewares, middlewares...)
}

func (app *httpApplication) Get(path string, handler HandleFunc, middlewares ...Middleware) {
	app.mux.HandleFunc(http.MethodGet+" "+path, func(w http.ResponseWriter, r *http.Request) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newMuxContext(w, setParam(path, r), &app.cfg.KafkaConfig, app.log))
	})
}

func (app *httpApplication) Post(path string, handler HandleFunc, middlewares ...Middleware) {
	app.mux.HandleFunc(http.MethodPost+" "+path, func(w http.ResponseWriter, r *http.Request) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newMuxContext(w, r, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *httpApplication) Put(path string, handler HandleFunc, middlewares ...Middleware) {
	app.mux.HandleFunc(http.MethodPut+" "+path, func(w http.ResponseWriter, r *http.Request) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newMuxContext(w, r, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *httpApplication) Delete(path string, handler HandleFunc, middlewares ...Middleware) {
	app.mux.HandleFunc(http.MethodDelete+" "+path, func(w http.ResponseWriter, r *http.Request) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newMuxContext(w, r, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *httpApplication) Patch(path string, handler HandleFunc, middlewares ...Middleware) {
	app.mux.HandleFunc(http.MethodPatch+" "+path, func(w http.ResponseWriter, r *http.Request) {
		preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newMuxContext(w, r, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *httpApplication) Register() *http.Server {
	server := http.Server{
		Handler:      app.mux,
		Addr:         ":" + app.cfg.AppConfig.Port,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
	}

	return &server

	// shutdown := make(chan error)

	// go func() {
	// 	quit := make(chan os.Signal, 1)

	// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 	<-quit

	// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// 	defer cancel()

	// 	log.Printf("Shutdown server: %s\n", server.Addr)
	// 	shutdown <- server.Shutdown(ctx)
	// }()

	// log.Printf("Start server: %s\n", server.Addr)
	// err := server.ListenAndServe()
	// if !errors.Is(err, http.ErrServerClosed) {
	// 	shutdown <- err
	// 	log.Fatal(err)
	// }

	// err = <-shutdown
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Server gracefully stopped")
}
