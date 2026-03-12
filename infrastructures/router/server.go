package router

import (
	stdContext "context"
	"fmt"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sirupsen/logrus"

	// For Swagger
	_ "github.com/mobigen/golang-web-template/docs/swagger"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	// and handles the control to the centralized HTTPErrorHandler.
	r.Use(middleware.Recover())

	// CORS default
	// Echo v5: AllowOrigins 또는 UnsafeAllowOriginFunc 필수. 모든 원격지에서 오는 모든 메서드를 허용합니다.
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

	// Swagger
	r.GET("/swagger/*", echoSwagger.WrapHandler)

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
