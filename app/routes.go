package app

import (
	"fmt"
	"net/http"

	"github.com/miodzie/gobroke/notifying"
)

// STOP THINKING ABOUT THE DESIGN
// JUST GET IT WORKING FIRST YOU GIT
// ONCE ITS WORKING, YOU CAN STRUCTURE AFTER
// THE DATABASE IS A DETAIL

// One place for all routes
func (s *Server) routes() {
	// TODO: Add routes
	s.router.HandlerFunc("GET", "/ping", s.handlePing())

	// TODO: implement
	s.router.HandlerFunc("GET", "/triggers", s.handleTriggersGet(s.trigRepo))
	s.router.HandlerFunc("POST", "/triggers", s.handleTriggersCreate(s.trigRepo))
}

func (s *Server) handleTriggersGet(repo notifying.TriggerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trigs, err := repo.GetAll()
		if err != nil {
			s.error(w, r, err)
			return
		}

		s.response(w, r, trigs, http.StatusOK)
	}
}

func (s *Server) handleTriggersCreate(repo notifying.TriggerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			Symbol    string `json:"symbol"`
			Threshold string `json:"threshold"`
			Disabled  bool   `json:"disabled"`
		}{}

		// Parse request.
		err := s.decode(w, r, req)
		if err != nil {
			s.error(w, r, err)
			return
		}

		// TODO: Validate.

		// Store.
		fmt.Printf("%v\n", req)
		trig := &notifying.Trigger{
			Symbol:    req.Symbol,
			Threshold: req.Threshold,
			Disabled:  req.Disabled,
		}

		repo.Save(trig)

		s.response(w, r, trig, http.StatusCreated)
	}
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Message string `json:"message"`
		}{"pong"}
		s.response(w, r, data, 200)
	}
}
