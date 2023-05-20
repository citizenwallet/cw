package router

import (
	"fmt"
	"net/http"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/response"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/community"
	"github.com/daobrussels/cw/pkg/hello"
	"github.com/daobrussels/cw/pkg/push"
	"github.com/daobrussels/cw/pkg/server"
	"github.com/daobrussels/cw/pkg/token"
	"github.com/daobrussels/cw/pkg/transaction"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	s  *supply.Supply
	es *ethrequest.EthService
	c  *community.Community
}

func NewServer(s *supply.Supply,
	es *ethrequest.EthService,
	c *community.Community) server.Server {
	return &Router{
		s,
		es,
		c,
	}
}

// implement the Server interface
func (r *Router) Start(port int) error {
	responder := response.NewResponder(r.s)

	cr := chi.NewRouter()

	// configure middleware
	cr.Use(OptionsMiddleware)
	cr.Use(HealthMiddleware)
	cr.Use(middleware.Compress(9))
	cr.Use(createSignatureMiddleware(r.s.PrivateHexKey))

	// instantiate handlers
	hello := hello.NewHandlers(r.c.Chain, responder)
	transaction := transaction.NewHandlers(&r.c.Chain, r.s, r.es)
	community := community.NewHandlers(responder, r.c)
	token := token.NewHandlers()
	push := push.NewHandlers()

	// configure routes
	cr.Get("/hello", hello.Hello)

	cr.Post("/transaction", transaction.Send)

	cr.Route("/gateway", func(cr chi.Router) {
		cr.Get("/", community.Config)
	})

	cr.Route("/token", func(cr chi.Router) {
		cr.Post("/mint", token.Mint)
		cr.Post("/burn", token.Burn)
	})

	cr.Route("/push", func(cr chi.Router) {
		cr.Put("/associate", push.Associate)
		cr.Delete("/dissociate", push.Dissociate)
	})

	// start the server
	return http.ListenAndServe(fmt.Sprintf(":%v", port), cr)
}
