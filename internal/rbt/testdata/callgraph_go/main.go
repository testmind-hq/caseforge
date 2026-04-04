// internal/rbt/testdata/callgraph_go/main.go
package main

import (
	"net/http"
	"testapp/handler"
	"testapp/repo"
)

func main() {
	mux := http.NewServeMux()
	r := &repo.MySQLRepo{} // RTA sees this instantiation
	handler.Register(mux, r)
	_ = http.ListenAndServe(":8080", mux)
}
