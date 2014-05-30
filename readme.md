Middleware
----------


This is as bare bones as possible, a reimagining of the middleware interface compatible with negroni built on top of the standard library container/list. It is written as an extended version of the net/http standard library, with the hope that it could be added to the standard library and people can come to some agreement on the interface to write middleware. I would rather have a standard than multiple solutions even if there are benefits to the differing implementations.