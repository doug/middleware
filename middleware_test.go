package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func stringwrap(base *string, pre, post string) Middleware {
	return MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.Handler) {
		*base += pre
		next.ServeHTTP(rw, r)
		*base += post
	})
}

func TestMiddlewareServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	s := NewStack()
	s.Use(stringwrap(&result, "a", "g"))
	mark := s.Use(stringwrap(&result, "b", "f"))
	s.Use(stringwrap(&result, "d", ""))

	// Insert at point
	s.InsertAfter(stringwrap(&result, "c", "e"), mark)

	// Append a whole stack
	s2 := NewStack()
	s2.Use(stringwrap(&result, "alpha_", "_omega"))
	s.PushFrontList(s2.List)

	request, _ := http.NewRequest("GET", "http://foobar.com", nil)
	s.ServeHTTP(response, request)

	expected := "alpha_abcdefg_omega"
	if result != expected {
		t.Errorf("Invalid result, expected %s got %s\n", expected, result)
	}
}

func TestMiddlewareCompose(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	m1 := Compose(stringwrap(&result, "a", "f"))
	m2 := Compose(stringwrap(&result, "b", "e"))
	m3 := Compose(stringwrap(&result, "c", "d"))

	s := m1(m2(m3(http.DefaultServeMux)))

	request, _ := http.NewRequest("GET", "http://foobar.com", nil)
	s.ServeHTTP(response, request)

	expected := "abcdef"
	if result != expected {
		t.Errorf("Invalid result, expected %s got %s\n", expected, result)
	}
}
