package routes

import (
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type routeMethods map[string]http.Handler

type Router struct {
	mux         *http.ServeMux
	prefix      string
	middlewares []Middleware // middlewares applied to handlers registered on this Router
	routes      map[string]routeMethods
}

// NewRouter creates a top-level Router.
func NewRouter() *Router {
	return &Router{
		mux:    http.NewServeMux(),
		routes: make(map[string]routeMethods),
	}
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Group(prefix string) *Router {
	p := strings.TrimSuffix(prefix, "/")
	// copy middlewares to preserve parent middleware chain
	mws := make([]Middleware, len(r.middlewares))
	copy(mws, r.middlewares)

	return &Router{
		mux:         r.mux,
		prefix:      r.prefix + p,
		middlewares: mws,
		routes:      r.routes,
	}
}

// full pattern = router.prefix + pattern
func (r *Router) fullPattern(pattern string) string {
	if pattern == "" {
		return r.prefix
	}
	// ensure leading slash on pattern
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	return r.prefix + pattern
}

func (r *Router) Handle(method, pattern string, handler http.Handler) {
	p := r.fullPattern(pattern)

	// initialize if path not registered
	if r.routes[p] == nil {
		r.routes[p] = make(routeMethods)

		r.mux.Handle(p, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			m := r.routes[p][req.Method]
			if m == nil {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			m.ServeHTTP(w, req)
		}))
	}

	// apply middlewares
	h := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	// register method-specific handler inside the map
	r.routes[p][method] = h
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
