# Middleware

Simple middleware on top of container/list. Middleware interface is
`func(rw http.ResponseWriter, r *http.Response, next http.Handler)`.

## Example Use

~~~ go

router := mux.NewRouter()
router.HandleFunc("/", SomeHandler)

s := middleware.NewStack()
s.Use(Middleware1)
element2 := s.Use(Middleware2)
s.Use(Middleware4)
s.UseHandler(router)

// You can modify the stack because it is a container/list
s.InsertAfter(Middleware3, element2)

// Or compose some set of middleware and use it
s2 := middleware.NewStack()
s2.Use(MiddlewareA)
s2.Use(MiddlewareB)
// Add this middleware stack to the front of the other one.
s.PushFrontList(s2.List)

// Compose converts a Middleware into a func(http.Handler)http.Handler
// so it can be called with Alice or just composing(functions(like(this))).
m1 := middleware.Compose(MiddlewareA)
m2 := middleware.Compose(MiddlewareB)
m3 := middleware.Compose(MiddlewareC)

// Stack 3
s3 := m1(m2(m3(http.DefaultServeMux)))

~~~~
