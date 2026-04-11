package engine

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Middleware func(http.Handler) http.Handler

type Route struct {
	Method      string
	Endpoint    string
	Description string
	Controller  HandlerFunc
	Middleware  []Middleware
}

func (e *HttpEngine) AddRoutes(routes []Route) {
	for _, route := range routes {
		var handler http.Handler = http.HandlerFunc(route.Controller)

		for i := len(route.Middleware) - 1; i >= 0; i-- {
			handler = route.Middleware[i](handler)
		}

		handlerCopy := handler

		finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := e.SetContext(r.Context(), nil)
			e.Log.HTTP(ctx, r.Method, r.Pattern)
			handlerCopy.ServeHTTP(w, r.WithContext(ctx))
		})

		switch route.Method {
		case http.MethodGet:
			e.router.Get(route.Endpoint, finalHandler)
		case http.MethodPost:
			e.router.Post(route.Endpoint, finalHandler)
		case http.MethodPut:
			e.router.Put(route.Endpoint, finalHandler)
		case http.MethodPatch:
			e.router.Patch(route.Endpoint, finalHandler)
		case http.MethodDelete:
			e.router.Delete(route.Endpoint, finalHandler)
		default:
			panic("unsupported method: " + route.Method)
		}

		fmt.Printf("%s : %s - %s \n", route.Method, route.Endpoint, route.Description)
	}
}
