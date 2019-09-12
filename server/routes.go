package server

import (
	"net/http"
	"strconv"

	"github.com/LucaCtt/thelist/common/store"
	"github.com/LucaCtt/thelist/errors"
	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.HandleFunc("/api", s.useStripSlashes(s.useRecoverer(s.useContentType("application/json", s.handleAPI()))))
}

func getID(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return 0, errors.E("invalid value for id param", err, errors.CodeBadValue)
	}

	return uint(id), nil
}

func (s *Server) handleAPI() http.HandlerFunc {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		items, err := s.store.All()
		if err != nil {
			s.error(w, r, errors.E("server get items failed", err))
			return
		}

		s.respond(w, r, items, http.StatusOK)
	})

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			s.error(w, r, errors.E("error parsing id param", err))
			return
		}

		show, err := s.store.Get(id)
		if err != nil {
			if errors.Code(err) == errors.CodeNotFound {
				s.error(w, r, errors.E("show not found", err))
			}
			s.error(w, r, errors.E("server get item failed", err))
			return
		}

		s.respond(w, r, show, http.StatusOK)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var show store.Item
		if err := s.decode(w, r, &show); err != nil {
			s.error(w, r, errors.E("invalid value for body", err, errors.CodeBadValue))
			return
		}

		if err := s.store.Create(&show); err != nil {
			s.error(w, r, errors.E("server create item failed", err))
			return
		}

		s.respond(w, r, nil, http.StatusNoContent)
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			s.error(w, r, errors.E("error parsing id param", err))
			return
		}

		if err := s.store.Delete(id); err != nil {
			s.error(w, r, errors.E("server delete item failed", err))
			return
		}

		s.respond(w, r, nil, http.StatusNoContent)
	})

	return router.ServeHTTP
}
