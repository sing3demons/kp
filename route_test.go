package kp

import (
	"syscall"
	"testing"
	"time"
)

func TestHttpApplicationStart(t *testing.T) {
	cfg := Config{
		AppConfig: AppConfig{
			Port: "8080",
		},
	}

	app := newServer(cfg).(*httpApplication)

	go func() {
		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	app.Start()
}
