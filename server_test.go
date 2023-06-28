package main

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func TestDarkMode(t *testing.T) {
	store := &MockStore{}
	server := NewServer(store)

	t.Run("get page with dark mode cookie set", func(t *testing.T) {
		cookie := http.Cookie{Name: "dark-mode", Value: "enabled"}
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&cookie)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), `<input id="toggle-btn" checked type="checkbox" class="peer opacity-0 w-0 h-0">`)
	})

	t.Run("get page with dark mode cookie not set", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), `<input id="toggle-btn"  type="checkbox" class="peer opacity-0 w-0 h-0">`)
	})
}

func TestVersion(t *testing.T) {
	store := &MockStore{}
	server := NewServer(store)

	t.Run("get page with version set", func(t *testing.T) {
		request := newIndexRequest()
		response := httptest.NewRecorder()

		expectedVersion := os.Getenv("VERSION")

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), expectedVersion)
	})
}
