package middleware

import (
	"container/list"
	"net/http"
)

// Middleware handler is an interface that objects can implement to be registered to serve as middleware
// in the stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next MiddlewareFunc.
// passed in.
type Middleware interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// MiddlewareFunc is an adapter to allow the use of ordinary functions as middleware handlers.
// If f is a function with the appropriate signature, MiddlewareFunc(f) is a Middleware object that calls f.
type MiddlewareFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

// Wrap converts a Handler into a Middleware so it can be used as a
// middleware. The next HandlerFunc is automatically called after the Middleware
// is executed.
func Wrap(handler http.HandlerFunc) Middleware {
	return MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

// Compose converts a Middleware into a func(http.HandlerFunc)http.HandlerFunc
// so it can be called with Alice or just composing(functions(like(this))).
func Compose(m Middleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			m.ServeHTTP(rw, r, next.ServeHTTP)
		})
	}
}

type middleware list.Element

func (m *middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	e := (*list.Element)(m)
	next := (*middleware)(e.Next())
	h := e.Value.(Middleware)
	if next == nil {
		h.ServeHTTP(rw, r, voidHandler)
		return
	}
	h.ServeHTTP(rw, r, next.ServeHTTP)
}

// Stack is a linked list stack of middleware
type Stack struct {
	*list.List
}

// NewStack returns a new linked list Stack of middlware
func NewStack() *Stack {
	return &Stack{list.New()}
}

func (s *Stack) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	front := (*middleware)(s.Front())
	if front != nil {
		front.ServeHTTP(rw, r)
	}
}

// Use adds a Middleware onto the middleware stack. Middlewares are invoked in the order they are added unless otherwise specified.
func (s *Stack) Use(handler Middleware) *list.Element {
	return s.PushBack(handler)
}

// UseHandler adds a Handler onto the middleware stack. Handlers are invoked in the order they are added unless otherwise specified.
func (s *Stack) UseHandler(handler http.HandlerFunc) *list.Element {
	return s.Use(Wrap(handler))
}

var voidHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
