// internal/rbt/testdata/callgraph_go/handler/handler.go
package handler

import (
	"net/http"
	"testapp/repo"
	"testapp/service"
)

// Register wires the /users route. This is the route-registering file.
func Register(mux *http.ServeMux, r repo.UserRepo) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		CreateUser(w, req, r)
	})
}

// CreateUser calls service.Process with the concrete repo.
func CreateUser(w http.ResponseWriter, r *http.Request, ur repo.UserRepo) {
	service.Process(ur)
}
