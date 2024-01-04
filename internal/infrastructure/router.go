package infrastructure

import (
	"errors"
	"net/http"

	"judge-opinioner/infrastructure/server"
	"judge-opinioner/ui"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	echo       *echo.Echo
	controller ui.Controller
}

func NewRouter(controller ui.Controller) *Router {
	e := echo.New()
	e.Use(corsMiddleware())

	server.RegisterHandlers(e, controller)

	return &Router{
		echo:       e,
		controller: controller,
	}
}

func (r *Router) Start(address string) error {
	if err := r.echo.Start(address); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func corsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
		},
	})
}
