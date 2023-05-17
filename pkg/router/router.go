package router

import (
	"fmt"
	"net/http"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/response"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/community"
	"github.com/daobrussels/cw/pkg/config"
	"github.com/daobrussels/cw/pkg/hello"
	"github.com/daobrussels/cw/pkg/push"
	"github.com/daobrussels/cw/pkg/server"
	"github.com/daobrussels/cw/pkg/token"
	"github.com/daobrussels/cw/pkg/transaction"
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

	s, err := supply.New(r.conf.SupplyWalletKey)
	if err != nil {
		return err
	}

	responder := response.NewResponder(s)
	ethservice, err := ethrequest.NewEthService(r.conf.Chain.RPC[0])
	if err != nil {
		return err
	}
	defer ethservice.Close()

	cr := chi.NewRouter()

	// configure middleware
	cr.Use(OptionsMiddleware)
	cr.Use(HealthMiddleware)
	cr.Use(middleware.Compress(9))
	cr.Use(createSignatureMiddleware(r.conf.SupplyWalletKey))

	// instantiate handlers
	hello := hello.NewHandlers(r.conf.Chain, responder)
	transaction := transaction.NewHandlers(&r.conf.Chain, s, ethservice)
	community := community.NewHandlers(s, ethservice, responder, &r.conf.Chain)
	token := token.NewHandlers()
	push := push.NewHandlers()

	// configure routes
	cr.Get("/hello", hello.Hello)

	cr.Post("/transaction", transaction.Send)

	cr.Route("/gateway", func(cr chi.Router) {
		cr.Post("/deploy", community.Deploy)
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
