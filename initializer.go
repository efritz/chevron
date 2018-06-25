package chevron

import (
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/process"
)

type (
	// ServerInitializer implements process.HTTPServerInitializer.
	ServerInitializer struct {
		Container   *nacelle.ServiceContainer `service:"container"`
		Logger      nacelle.Logger            `service:"logger"`
		initializer RouteInitializer
	}

	// RouteInitializer initializes a Router instance.
	RouteInitializer interface {
		// Init registers resources and middleware to the router.
		Init(config nacelle.Config, router Router) error
	}

	// RouteInitializerFunc is a function conforming to the RouteInitializer
	// interface.
	RouteInitializerFunc func(config nacelle.Config, router Router) error
)

// Init calls the wrapped function.
func (f RouteInitializerFunc) Init(config nacelle.Config, router Router) error {
	return f(config, router)
}

// NewInitializer creates a new ServerInitializer.
func NewInitializer(initializer RouteInitializer) process.HTTPServerInitializer {
	return &ServerInitializer{
		initializer: initializer,
	}
}

// Init creates a router which becomes the server's handler and calls the
// attached route initializer.
func (i *ServerInitializer) Init(config nacelle.Config, server *http.Server) error {
	// TODO - control additional configs with env vars
	router := NewRouter(i.Container, WithLogger(i.Logger))
	server.Handler = router
	return i.initializer.Init(config, router)
}