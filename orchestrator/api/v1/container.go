package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afif1400/ucmdls/orchestrator/docker"
	"github.com/afif1400/ucmdls/orchestrator/helpers"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	nats "github.com/docker/go-connections/nat"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

func HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}
}

// HandleContainers function to list all the containers
func HandleContainers(d *docker.Docker, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containers, err := d.Client.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

		// convert to json and write to response
		err = json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusOK, "success", containers))
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

	}
}

// HandleContainer function to get a container
func HandleContainer(d *docker.Docker, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containerId := mux.Vars(r)["id"]
		containerInspect, err := d.Client.ContainerInspect(ctx, containerId)
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

		// convert to json and write to response
		err = json.NewEncoder(w).Encode(containerInspect)
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

	}
}

// HandleContainerCreate function to create a container
func HandleContainerCreate(d *docker.Docker, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newport, err := nats.NewPort("tcp", "4200")
		// get container name from request body
		var containerRequest struct {
			Name string `json:"name"`
		}

		err = json.NewDecoder(r.Body).Decode(&containerRequest)
		if containerRequest.Name == "" {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusBadRequest, "failure", "Container name is required"))
			return
		}
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

		hostConfig := &container.HostConfig{
			PortBindings: nats.PortMap{
				newport: []nats.PortBinding{
					{
						HostPort: "5001",
					},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
			LogConfig: container.LogConfig{
				Type:   "json-file",
				Config: map[string]string{},
			},
		}

		networkConfig := &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{},
		}

		exposedPorts := map[nats.Port]struct{}{
			newport: {},
		}

		envVars := []string{
			"SIAB_PORT=4200",
			"SIAB_PASSWORD=123123123",
			"SIAB_SUDO=true",
		}
		containerConfig := &container.Config{
			Image:        "sspreitzer/shellinabox:latest",
			Env:          envVars,
			ExposedPorts: exposedPorts,
			Hostname:     fmt.Sprint("shellinabox-hostname"),
		}

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := d.Client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerRequest.Name)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// convert to json and write to response
		err = json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusOK, "created container", resp))
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

		// start the container and print the output to the response
		err = d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

		json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusCreated, "started container", resp))
	}
}

// HandleContainerStart function to start a container
func HandleContainerStart(d *docker.Docker, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containerId := mux.Vars(r)["id"]
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := d.Client.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
			if err != nil {
				json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
				return
			}
		}()
		wg.Wait()

		// convert to json and write to response
		err := json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusOK, "started container", nil))
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

	}
}

// HandleContainerStop function to stop a container
func HandleContainerStop(d *docker.Docker, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containerId := mux.Vars(r)["id"]

		// wait for the container to stop
		stopCh := make(chan error)
		go func() {
			err := d.Client.ContainerStop(ctx, containerId, nil)
			stopCh <- err
		}()

		select {
		case err := <-stopCh:
			if err != nil {
				json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
				return
			}
		case <-time.After(time.Second * 10):
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusRequestTimeout, "failure", "timeout"))
			return
		}

		json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusOK, "stopped container", nil))
	}
}

// HandleContainerRemove function to delete a container
func HandleContainerRemove(d *docker.Docker, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containerId := mux.Vars(r)["id"]
		err := d.Client.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
		if err != nil {
			json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusInternalServerError, "failure", err.Error()))
			return
		}

		json.NewEncoder(w).Encode(helpers.NewApiResponse(http.StatusOK, "deleted container", nil))
	}
}
