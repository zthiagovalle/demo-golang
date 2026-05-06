package shared

import "github.com/labstack/echo/v5"

type Route struct {
	URI     string
	Method  string
	Handler echo.HandlerFunc
}

func RegisterRoutes(e *echo.Echo, routes []Route) {
	for _, r := range routes {
		e.Add(r.Method, r.URI, r.Handler)
	}
}
