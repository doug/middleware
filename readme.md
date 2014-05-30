# Middleware

This is as bare bones as possible, a reimagining of the middleware interface compatible with negroni built on top of the standard library container/list. It is written with the intention of integration with the net/http standard library, so names etc are such that they do not conflict. I would hope that some day in the future it or something like it could be added to the standard library and people can come to some agreement on the interface to write middleware. I would rather have a standard than multiple solutions even if there are benefits to the differing implementations.

I feel strongly that any middleware should have an interface only using the net/http standard. The reason I chose `func(rw http.ResponseWriter, r *http.Response, next http.HandlerFunc)` over the other most common approach of just `func(http.Handler) http.Handler` is because with a linked list you can make partial stacks and compose them with other ones, so someone can precombine a set of middleware to import and use as a unit. Additionally, the code produces by calling `next(rw,r)` is simpler than a returned closure function that the other approach generally uses. 

As far as needing some sort of context for exchanging information between middleware layers. I leave that like the choice of mux undecided, there are plenty out there, or you could just append some information to the request header. Personally, I use gorilla for both. http://www.gorillatoolkit.org/pkg/mux and http://www.gorillatoolkit.org/pkg/context.

## Example Use

~~~ go

// The middleware interface is compatible with negroni
// func(rw http.ResponseWriter, r *http.Response, next http.HandlerFunc)
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

~~~~

### Blogs about Go middleware
  - http://stephensearles.com/?p=254
  - http://codegangsta.io/blog/2014/05/19/my-thoughts-on-martini/
  - http://justinas.org/writing-http-middleware-in-go/
  - http://justinas.org/alice-painless-middleware-chaining-for-go/
  - http://www.reddit.com/r/golang/comments/252wjh/are_you_using_golang_for_webapi_development_what/

### Other middleware
  - https://github.com/gorilla/handlers
  - https://github.com/stephens2424/muxchain
  - https://github.com/go-martini/martini
  - https://github.com/codegangsta/negroni
  - https://github.com/justinas/alice

### Web frameworks with middleware
  - https://github.com/yosssi/galaxy
  - http://goji.io
