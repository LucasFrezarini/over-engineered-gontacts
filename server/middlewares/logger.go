package middlewares

import (
	"fmt"
	"net/http/httputil"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Container struct {
	logger *zap.Logger
}

func ProvideMiddlewaresContainer(logger *zap.Logger) *Container {
	return &Container{logger: logger.Named("HTTPLogger")}
}

func (ct *Container) ZapHTTPLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dump, err := httputil.DumpRequest(c.Request(), true)
		if err != nil {
			ct.logger.Error(fmt.Sprintf("error while dumping request: %v", err))
		}

		ct.logger.Info(fmt.Sprintf("received incoming request: %s", dump))
		return next(c)
	}
}

var MiddlewaresSet = wire.NewSet(ProvideMiddlewaresContainer)
