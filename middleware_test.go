package http

import (
	. "net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestMiddlewareServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	s := NewStack()
	s.Use(MiddlewareFunc(func(rw ResponseWriter, r *Request, next HandlerFunc) {
		result += "foo"
		next(rw, r)
		result += "ban"
	}))
	s.Use(MiddlewareFunc(func(rw ResponseWriter, r *Request, next HandlerFunc) {
		result += "bar"
		next(rw, r)
		result += "baz"
	}))
	s.Use(MiddlewareFunc(func(rw ResponseWriter, r *Request, next HandlerFunc) {
		result += "bat"
		rw.WriteHeader(StatusBadRequest)
	}))

	s.ServeHTTP(response, (*Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, StatusBadRequest)
}
