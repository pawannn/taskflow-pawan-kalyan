package engine

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Middleware func(http.Handler) http.Handler

type Route struct {
	Method      string
	Path        string
	Description string
	Controller  HandlerFunc
	Middleware  []Middleware
}

func (e *HttpEngine) AddRoutes(routes []Route) {
	for _, route := range routes {
		var handler http.Handler = http.HandlerFunc(route.Controller)

		// middlewares are executed in reverse order
		for i := len(route.Middleware) - 1; i >= 0; i-- {
			handler = route.Middleware[i](handler)
		}

		finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method Not Allowed"))
				return
			}

			ctx := e.SetContext(r.Context(), nil)
			r = r.WithContext(ctx)

			handler.ServeHTTP(w, r)
		})

		fmt.Printf("%s : %s - %s \n", route.Method, route.Path, route.Description)
		e.mux.Handle(route.Path, finalHandler)
	}
}
