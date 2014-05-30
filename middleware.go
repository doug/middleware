package http

import (
	"container/list"
	. "net/http"
)

// Middleware handler is an interface that objects can implement to be registered to serve as middleware
// in the stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next MiddlewareFunc.
// passed in.
type Middleware interface {
	ServeHTTP(rw ResponseWriter, r *Request, next HandlerFunc)
}

// MiddlewareFunc is an adapter to allow the use of ordinary functions as middleware handlers.
// If f is a function with the appropriate signature, MiddlewareFunc(f) is a Middleware object that calls f.
type MiddlewareFunc func(rw ResponseWriter, r *Request, next HandlerFunc)

func (h MiddlewareFunc) ServeHTTP(rw ResponseWriter, r *Request, next HandlerFunc) {
	h(rw, r, next)
}

// Wrap converts a Handler into a Middleware so it can be used as a
// middleware. The next HandlerFunc is automatically called after the Middleware
// is executed.
func Wrap(handler Handler) Middleware {
	return MiddlewareFunc(func(rw ResponseWriter, r *Request, next HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

type middleware list.Element

func (m *middleware) ServeHTTP(rw ResponseWriter, r *Request) {
	e := (*list.Element)(m)
	next := (*middleware)(e.Next())
	h := e.Value.(Middleware)
	if next == nil {
		h.ServeHTTP(rw, r, voidMiddleware)
		return
	}
	h.ServeHTTP(rw, r, next.ServeHTTP)
}

// Stack is a linked list stack of middleware
type stack struct {
	*list.List
}

// New returns a new linked list stack of middlware
func NewStack() *stack {
	return &stack{list.New()}
}

func (s *stack) ServeHTTP(rw ResponseWriter, r *Request) {
	front := (*middleware)(s.Front())
	if front != nil {
		front.ServeHTTP(rw, r)
	}
}

// Use adds a Middleware onto the middleware stack. Middlewares are invoked in the order they are added unless otherwise specified.
func (s *stack) Use(handler Middleware) {
	s.PushBack(handler)
}

// UseHandler adds a Handler onto the middleware stack. Handlers are invoked in the order they are added unless otherwise specified.
func (s *stack) UseHandler(handler Handler) {
	s.Use(Wrap(handler))
}

func voidMiddleware(rw ResponseWriter, r *Request) {
	// do nothing
}
