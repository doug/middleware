package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func stringwrap(base *string, pre, post string) Middleware {
	return MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		*base += pre
		next(rw, r)
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

	s.ServeHTTP(response, (*http.Request)(nil))

	expected := "alpha_abcdefg_omega"
	if result != expected {
		t.Errorf("Invalid result, expected %s got %s\n", expected, result)
	}
}
