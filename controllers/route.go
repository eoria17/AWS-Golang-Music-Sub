package controllers

import (
	"github.com/eoria17/AWS-Golang-Music-Sub/controllers/login"
	"github.com/gorilla/mux"
)

func Route(r *mux.Router) {
	r.HandleFunc("/", login.Login)
}
