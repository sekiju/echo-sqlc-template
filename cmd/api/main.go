package main

import (
	"context"
	"echo-sqlc-template/internal/config"
	"echo-sqlc-template/internal/controller"
	"echo-sqlc-template/internal/database"
	"echo-sqlc-template/internal/services/mail"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	err := config.Load()
	if err != nil {
		return err
	}

	err = database.Connect()
	if err != nil {
		return fmt.Errorf("error intializing database: %v", err)
	}

	mail.Connect()

	e := echo.New()
	e.Debug = config.Data.Application.Development
	controller.Create(e)

	go func() {
		if err := e.Start(config.Data.Application.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}
