package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"slow/config"
	"slow/controller/petrol"
	"slow/dal"
	"slow/logger"
	"slow/util/custom_validator"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
)

func main() {
	err := config.Initialize()
	if err != nil {
		logger.Get().Error().Err(err).Msg("failed to load config")
		panic(err)
	}
	err = dal.Initialize()
	if err != nil {
		logger.Get().Error().Err(err).Msg("failed to connect to database")
		panic(err)
	}
	app := &cli.App{
		Name:  "slow",
		Usage: "slow is fast.",
		Action: func(ctx *cli.Context) error {
			startServer()
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Get().Error().Err(err).Msg("failed to start the application")
		panic(err)
	}
}

func startServer() {
	e := echo.New()
	e.HideBanner = true
	e.Validator = custom_validator.NewCustomValidator()
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogLatency:   true,
		LogRemoteIP:  true,
		LogMethod:    true,
		LogURI:       true,
		LogRequestID: true,
		LogStatus:    true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Get().Info().
					Str("remote_ip", v.RemoteIP).
					Str("method", v.Method).
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("request_id", v.RequestID).
					Msg("successful request")
			} else {
				logger.Get().Error().
					Str("remote_ip", v.RemoteIP).
					Str("method", v.Method).
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("request_id", v.RequestID).
					Err(v.Error).
					Msg("failed request")
			}
			return v.Error
		},
	}))
	e.Use(middleware.Recover())
	petrol.RegisterRoutes(e.Group("/petrol"))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Get().Fatal().Err(err).Msg("failed server start")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Get().Fatal().Err(err).Msg("failed server shutdown")
	}
}
