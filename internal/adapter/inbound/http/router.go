package http

import (
	stdContext "context"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/dto"
	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/handler"
	"github.com/mobigen/golang-web-template/internal/domain"

	// For Swagger
	_ "github.com/mobigen/golang-web-template/docs/swagger"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// Router echo.Echo
type Router struct {
	*echo.Echo
	Debug        bool
	ctx          stdContext.Context
	serverCancel stdContext.CancelFunc
}

// Init Echo Framework Initialize
func Init(log *logrus.Logger, debug bool) (r *Router, err error) {
	r = &Router{Echo: echo.New(), Debug: debug}
	ctx, cancel := stdContext.WithCancel(stdContext.Background())
	r.serverCancel = cancel
	r.ctx = ctx

	// Recover returns a middleware which recovers from panics anywhere in the chain
	r.Use(middleware.Recover())

	// CORS default
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Request Logger
	r.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper:    r.LoggerSkipper,
		LogLatency: true,
		LogMethod:  true,
		LogURI:     true,
		LogStatus:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			log.Infof("%s [DEBU] [echo-framework   :  - ] [ Router ] %s %s %d Latency[ %s ]",
				v.StartTime.Format("2006-01-02 15:04:05.000"),
				v.Method,
				v.URI,
				v.Status,
				v.Latency.Round(time.Millisecond),
			)
			return nil
		},
	}))

	// 전역 에러 핸들러
	r.HTTPErrorHandler = func(c *echo.Context, err error) {
		var ae domain.AppError
		if errors.As(err, &ae) {
			log.Warnf("[HTTPErrorHandler] business error: %v", err)
			_ = handler.FailApp(c, ae)
			return
		}
		if he, ok := err.(*echo.HTTPError); ok {
			_ = handler.Fail(c, he.Code, dto.ErrRouteNotFound, fmt.Sprintf("%v", he.Message))
			return
		}
		log.Errorf("[HTTPErrorHandler] unexpected error: %v", err)
		_ = handler.Fail(c, 500, dto.ErrInternalServer, "")
	}

	// Swagger
	r.GET("/swagger/*", echoSwagger.WrapHandlerV3)

	return r, nil
}

// EnableDebug debug mode on
func (r *Router) EnableDebug() {
	r.Debug = true
}

// DisableDebug disable debug
func (r *Router) DisableDebug() {
	r.Debug = false
}

// LoggerSkipper .. logger skipper
func (r *Router) LoggerSkipper(c *echo.Context) bool {
	if r.Debug {
		return false
	}
	return true
}

// Run echo framework
func (r *Router) Run(listenAddr string) error {
	if r == nil {
		return fmt.Errorf("ERROR. Router Not Initialize")
	}
	sc := echo.StartConfig{
		Address:    listenAddr,
		HideBanner: true,
		HidePort:   true,
	}
	return sc.Start(r.ctx, r.Echo)
}

// Shutdown echo framework
func (r *Router) Shutdown() error {
	if r.serverCancel != nil {
		r.serverCancel()
	}
	return nil
}
