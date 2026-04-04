package handler

import "net/http"

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", CreateUser)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	validate(r)
}
