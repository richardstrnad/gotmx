package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockStore struct{}

func (m *MockStore) GetTask(id int) (Task, error) {
	return Task{}, nil
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("response body is wrong, got %q should contain %q", got, want)
	}
}

func newIndexRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	return req
}

func TestIndexPage(t *testing.T) {
	store := &MockStore{}
	server := NewServer(store)

	t.Run("Get full page", func(t *testing.T) {
		request := newIndexRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "index")
	})

	t.Run("Get invalid page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/invalid678031324", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("Verify that the title is properly set", func(t *testing.T) {
		request := newIndexRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "<title>Index</title>")
	})
}
