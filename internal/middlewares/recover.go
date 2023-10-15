package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func Recover(log *zap.Logger) echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 4 << 10, // 4 KB
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error(
				"recover from panic",
				zap.Error(err),
				zap.ByteString("stack", stack),
			)
			return nil
		},
	})
}
