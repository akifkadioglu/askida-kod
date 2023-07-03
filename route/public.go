package route

import (
	"github.com/akifkadioglu/askida-kod/pkg/home"
	"github.com/go-chi/chi/v5"
)

func public(r chi.Router) {

	r.Group(func(r chi.Router) {
		r.Route("/public", func(r chi.Router) {
			r.Get("/", home.Home)
		})
	})
}
