package main

import (
	"context"
	apiv1 "github.com/afif1400/ucmdls/orchestrator/api/v1"
	"github.com/afif1400/ucmdls/orchestrator/docker"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	router *mux.Router
	docker *docker.Docker
}

func NewServer() *Server {

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		panic("error initializing docker client")
	}

	return &Server{
		router: mux.NewRouter(),
		docker: dockerClient,
	}
}

func (s *Server) Routes() *mux.Router {
	router := s.router
	ctx := context.Background()

	router.HandleFunc("/api/orchestrator", apiv1.HandleIndex()).Methods("GET")
	router.HandleFunc("/api/orchestrator/containers", apiv1.HandleContainers(s.docker, ctx)).Methods("GET")
	router.HandleFunc("/api/orchestrator/containers", apiv1.HandleContainerCreate(s.docker, ctx)).Methods("POST")
	router.HandleFunc("/api/orchestrator/containers/{id}", apiv1.HandleContainer(s.docker, ctx)).Methods("GET")
	router.HandleFunc("/api/orchestrator/containers/{id}/remove", apiv1.HandleContainerRemove(s.docker, ctx)).Methods("DELETE")
	router.HandleFunc("/api/orchestrator/containers/{id}/start", apiv1.HandleContainerStart(s.docker, ctx)).Methods("POST")
	router.HandleFunc("/api/orchestrator/containers/{id}/stop", apiv1.HandleContainerStop(s.docker, ctx)).Methods("POST")
	router.HandleFunc("/api/orchestrator/containers/commit", apiv1.HandleContainerCommit(s.docker, ctx)).Methods("POST")
	//router.HandleFunc("/containers/{id}/restart", apiv1.HandleContainerRestart(s.docker)).Methods("POST")
	//router.HandleFunc("/containers/{id}/logs", apiv1.HandleContainerLogs(s.docker)).Methods("GET")
	//router.HandleFunc("/containers/{id}/exec", apiv1.HandleContainerExec(s.docker)).Methods("POST")
	//router.HandleFunc("/images", apiv1.HandleImages(s.docker)).Methods("GET")
	//router.HandleFunc("/images/{id}", apiv1.HandleImage(s.docker)).Methods("GET")
	//router.HandleFunc("/images/{id}/remove", apiv1.HandleImageRemove(s.docker)).Methods("POST")
	//router.HandleFunc("/images/{id}/push", apiv1.HandleImagePush(s.docker)).Methods("POST")
	//router.HandleFunc("/images/{id}/tag", apiv1.HandleImageTag(s.docker)).Methods("POST")
	//router.HandleFunc("/images/create", apiv1.HandleImageCreate(s.docker)).Methods("POST")
	//router.HandleFunc("/images/{id}/load", apiv1.HandleImageLoad(s.docker)).Methods("POST")
	//router.HandleFunc("/images/load", apiv1.HandleImageLoad(s.docker)).Methods("POST")
	//router.HandleFunc("/images/{id}/save", apiv1.HandleImageSave(s.docker)).Methods("POST")
	//router.HandleFunc("/images/save", apiv1.HandleImageSave(s.docker)).Methods("POST")

	return router
}

func (s *Server) Run(r *mux.Router) {
	if err := http.ListenAndServe(":5000", r); err != nil {
		panic(err)
	}
}
