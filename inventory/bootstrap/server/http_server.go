package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"synapsis/inventory/config"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	EchoInstance struct {
		Router *echo.Echo
		Config config.AppConfig
	}
)

func NewEchoInstance(r *echo.Echo, cfg config.AppConfig) *EchoInstance {
	return &EchoInstance{
		Router: r,
		Config: cfg,
	}
}

func (server *EchoInstance) runHttp() (err error) {
	port := fmt.Sprintf(":%s", server.Config.App.HttpPort)

	if err = server.Router.Start(port); err != nil {
		return err
	}

	return
}

func (server *EchoInstance) RunWithGracefullyShutdown() {
	// run group on another thread
	go func() {
		err := server.runHttp()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Router.Shutdown(ctx); err != nil {
		os.Exit(1)
	}
}
