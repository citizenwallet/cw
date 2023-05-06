package router

import (
	"fmt"
	"net/http"

	"github.com/daobrussels/cw/pkg/config"
	"github.com/daobrussels/cw/pkg/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	// ...
	conf *config.Config
}

func NewServer(conf *config.Config) server.Server {
	return &Router{
		conf: conf,
	}
}

// implement the Server interface
func (r *Router) Start(port int) error {

	cr := chi.NewRouter()

	// configure middleware
	cr.Use(middleware.Compress(5))
	cr.Use(HealthMiddleware)
	// TODO: signature middleware

	// configure routes

	// start the server
	return http.ListenAndServe(fmt.Sprintf(":%v", port), cr)
}
