# Middleware

This is as bare bones as possible, a reimagining of the middleware interface compatible with negroni built on top of the standard library container/list. It is written with the intention of integration with the net/http standard library, so names etc are such that they do not conflict. I would hope that some day in the future it or something like it could be added to the standard library and people can come to some agreement on the interface to write middleware. I would rather have a standard than multiple solutions even if there are benefits to the differing implementations.

## Example Use

~~~ go

// Similar to negroni in the interface 
// func(rw http.ResponseWriter, r *http.Response, next http.HandleFunc)
router := mux.NewRouter()
router.HandleFunc("/", SomeHandler)

s := middleware.NewStack()
s.Use(Middleware1)
element2 := s.Use(Middleware2)
s.Use(Middleware4)
s.UseHandler(router)

// You can modify the stack because it is a container/list
s.InsertAfter(Middleware3, element2)
s.InsertAfter(Middleware5, s.Get(Middleware4))
~~~~

### Other middleware
  - https://github.com/stephens2424/muxchain
  - https://github.com/go-martini/martini
  - https://github.com/codegangsta/negroni
  - https://github.com/justinas/alice

### web frameworks with middleware
  - https://github.com/yosssi/galaxy
  - http://goji.io
