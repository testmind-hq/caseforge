// internal/rbt/testdata/callgraph_go/handler/handler.go
package handler

import (
	"net/http"
	"testapp/middleware"
	"testapp/repo"
	"testapp/service"
)

// Register wires the /users route. This is the route-registering file.
func Register(mux *http.ServeMux, r repo.UserRepo) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		CreateUser(w, req, r)
	})
}

// CreateUser calls service.Process with the concrete repo and middleware.Process
// (so deep.go is two hops away from this route-registering file).
func CreateUser(w http.ResponseWriter, r *http.Request, ur repo.UserRepo) {
	service.Process(ur)
	middleware.Process()
}
