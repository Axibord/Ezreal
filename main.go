package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github/axibord/heimerdinger/internal/domain/middlewares"
	"github/axibord/heimerdinger/internal/web/pages"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	m "github.com/labstack/echo/v4/middleware"
)

func RenderTempl(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func main() {
	e := echo.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Static("/static", "static")
	e.Use(m.RequestID())
	e.Use(m.RequestLoggerWithConfig(m.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogLatency:  true,
		LogMethod:   true,
		HandleError: false, // forwards error to the global error handler
		LogValuesFunc: func(c echo.Context, v m.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQ",
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
					slog.String("uri", v.URI),
					slog.Float64("latency", v.Latency.Seconds()),
					slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQ_ERR",
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
					slog.String("uri", v.URI),
					slog.Float64("latency", v.Latency.Seconds()),
					slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
				)
			}
			return nil
		},
	}))

	// security middlewares
	e.Use(middlewares.CSPMiddleware())
	e.Use(m.Secure())
	e.Use(m.CORS())
	e.Use(m.CSRF())
	e.Use(m.Decompress())
	e.Use(m.GzipWithConfig(m.GzipConfig{
		Level: 5,
	}))

	// Routes
	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Build your dream app with Heimerdinger!",
		})
	})

	e.GET("/home", func(ctx echo.Context) error {
		return RenderTempl(ctx, http.StatusOK, pages.Index())
	})

	// not found 404
	e.GET("/*", func(ctx echo.Context) error {
		return RenderTempl(ctx, http.StatusNotFound, pages.NotFound())
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(":4000"); err != nil && err != http.ErrServerClosed {
			logger.LogAttrs(context.Background(), slog.LevelError, "ERROR", slog.String("err", err.Error()))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "ERROR", slog.String("err", err.Error()))
	}
}
