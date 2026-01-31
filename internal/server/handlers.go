package server

import (
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/email"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/imoveis"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/sliders"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/user"
)

// Handlers aggregates handler instances and shared services used by route registration.
type Handlers struct {
	User    *user.Handler
	Sliders *sliders.Handler
	Imoveis *imoveis.Handler
	Email   *email.Handler
}
