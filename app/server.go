package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/miodzie/gobroke/notifying"
	"github.com/miodzie/gobroke/pricing"
)

type Server struct {
	db           *sql.DB
	router       *httprouter.Router
	priceChecker *pricing.Checker
	storage      interface{}
	// TODO: Switch to a RepoManager to consolidate.
	trigRepo notifying.TriggerRepository
}

// NewServer creates a new Server instance.
// Don't setup dependencies within the new function.
func NewServer(db *sql.DB, checker *pricing.Checker, trigRepo notifying.TriggerRepository) (*Server, error) {
	s := &Server{
		db:           db,
		router:       httprouter.New(),
		priceChecker: checker,
		trigRepo:     trigRepo,
	}
	s.routes()

	return s, nil
}

// ServerHTTP satisfies the http.HTTPHandler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Response helper.
// Abstract responding and do the bare bones initially
// Later you can make this more sophisticated (if needed)
func (s *Server) response(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(data)
		return
	}
	w.WriteHeader(status)
}

// Decode helper.
// Abstract decoding and do the bare bones initially
// Later you can make this more sophisticated (if needed)
func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// error helper.
func (s *Server) error(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: better error handling.
	log.Err(err).Send()
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// ----NOTES----

// Return the handler
// Take arguments for handler-specific dependencies
// handleSendMagicLinkEmail(e EmailSender) http.HandlerFunc
// Too big? Have many servers
// HandlerFunc over Handler
func (s *Server) handleSomething(format string) http.HandlerFunc {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		fmt.Fprintf(w, format, "potato")
	}
}

// Middleware are just Go functions
// Wire Middleware up in routes.go
func (s *Server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if !currentUser(r).isAdmin {
		// 	http.NotFound(w, r)
		// 	return
		// }
		h(w, r)
	}
}

// Request and response data types
func (s *Server) handleGreet() http.HandlerFunc {
	type request struct {
		Name string
	}

	type response struct {
		Greeting string `json:"greeting"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// ...
	}
}
