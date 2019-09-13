package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LucaCtt/thelist/common"
	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/common/store"
	"github.com/LucaCtt/thelist/common/testutils"
	"github.com/LucaCtt/thelist/mocks"
)

func TestServer_handleAPIGetAll(t *testing.T) {
	tests := []struct {
		name   string
		items  []*store.Item
		movies []*client.Movie
		want   []*common.Show
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mocks.Store{}
			client := &mocks.Client{}

			store.On("All").Return(tt.items)
			client.On("GetMovie").Return(tt.movies)

			s := &Server{store, client, http.NewServeMux()}
			s.router.HandleFunc("/", s.handleAPIGetAll())

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			s.ServeHTTP(w, r)

			testutils.AssertEqual(t, w.Result().StatusCode, http.StatusOK)
		})
	}
}

func TestServer_handleAPISearch(t *testing.T) {
	type fields struct {
		store  store.Store
		client client.Client
		router *http.ServeMux
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				store:  tt.fields.store,
				client: tt.fields.client,
				router: tt.fields.router,
			}
			if got := s.handleAPISearch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.handleAPISearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_handleAPIDelete(t *testing.T) {
	type fields struct {
		store  store.Store
		client client.Client
		router *http.ServeMux
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				store:  tt.fields.store,
				client: tt.fields.client,
				router: tt.fields.router,
			}
			if got := s.handleAPIDelete(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.handleAPIDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_handleAPIPost(t *testing.T) {
	type fields struct {
		store  store.Store
		client client.Client
		router *http.ServeMux
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				store:  tt.fields.store,
				client: tt.fields.client,
				router: tt.fields.router,
			}
			if got := s.handleAPIPost(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.handleAPIPost() = %v, want %v", got, tt.want)
			}
		})
	}
}
