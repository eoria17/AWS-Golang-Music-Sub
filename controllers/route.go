package controllers

import (
	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AppEngine struct {
	Session     *session.Session
	Credentials *credentials.Credentials
}

func (ae AppEngine) Route(r *mux.Router) {
	r.HandleFunc("/", ae.Login)
}
