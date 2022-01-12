package main

import (
	"net/http"

	apiv1 "github.com/afif1400/ucmdls/manager/api/v1"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {

	return &Server{
		router: mux.NewRouter(),
	}
}

func (s *Server) Routes() *mux.Router {
	router := s.router

	router.HandleFunc("/api/manager/labs", apiv1.HandleLabs()).Methods("GET")

	return router
}

func (s *Server) Run(r *mux.Router) {
	if err := http.ListenAndServe(":5000", r); err != nil {
		panic(err)
	}
}
