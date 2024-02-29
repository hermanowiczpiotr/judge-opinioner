package infrastructure

import (
	"errors"
	"net/http"

	"judge-opinioner/internal/infrastructure/server"
	"judge-opinioner/internal/ui"

	internal_middlawere "judge-opinioner/internal/infrastructure/middleware"

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

	e.Use(internal_middlawere.JwtMiddleware)
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
