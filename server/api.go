package server

import (
	"net/http"
	"strconv"

	"github.com/LucaCtt/thelist/common"
	"github.com/LucaCtt/thelist/common/store"
	"github.com/LucaCtt/thelist/errors"
	"github.com/go-chi/chi"
)

func getID(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return 0, errors.E("invalid value for id param", err, errors.CodeBadValue)
	}

	return uint(id), nil
}

func (s *Server) handleAPIGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.store.All()
		if err != nil {
			s.error(w, r, errors.E("server get items failed", err))
			return
		}

		shows := make([]*common.Show, len(items))
		for i, item := range items {
			show, err := common.GetShow(s.client, item.ShowID, item.Type)
			if err != nil {
				s.error(w, r, errors.E("server get show info failed", err))
				return
			}
			shows[i] = show
		}

		s.respond(w, r, shows, http.StatusOK)
	}
}

func (s *Server) handleAPISearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) handleAPIDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func (s *Server) handleAPIPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
