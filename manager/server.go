package main

import (
	"net/http"

	apiv1 "github.com/afif1400/ucmdls/manager/api/v1"
	"github.com/afif1400/ucmdls/manager/database"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	router *mux.Router
	dbc    *mongo.Collection
}

func NewServer() *Server {
	// connet to the mongodb database
	collection := database.ConnectDB()

	return &Server{
		router: mux.NewRouter(),
		dbc:    collection,
	}
}

func (s *Server) Routes() *mux.Router {
	router := s.router

	router.HandleFunc("/api/manager/labs", apiv1.HandleLabs(s.dbc)).Methods("GET")
	router.HandleFunc("/api/manager/labs/{labID}", apiv1.HandleLab(s.dbc)).Methods("GET")
	router.HandleFunc("/api/manager/labs/create", apiv1.HandleLabCreate(s.dbc)).Methods("POST")
	router.HandleFunc("/api/manager/labs/{labID}/systems", apiv1.HandleLabAddSystems(s.dbc)).Methods("PUT")

	return router
}

func (s *Server) Run(r *mux.Router) {
	if err := http.ListenAndServe(":5000", r); err != nil {
		panic(err)
	}
}
