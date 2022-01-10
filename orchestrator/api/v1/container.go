package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afif1400/ucmdls/orchestrator/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	nats "github.com/docker/go-connections/nat"
	"github.com/gorilla/mux"
	"net/http"
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
			w.Write([]byte(err.Error()))
			return
		}

		// convert to json and write to response
		err = json.NewEncoder(w).Encode(containers)
		if err != nil {
			w.Write([]byte(err.Error()))
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
			w.Write([]byte(err.Error()))
			return
		}

		// convert to json and write to response
		err = json.NewEncoder(w).Encode(containerInspect)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

	}
}

// HandleContainerCreate function to create a container
func HandleContainerCreate(d *docker.Docker, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newport, err := nats.NewPort("tcp", "4200")
		containerName := "shellinabox-demo"
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

		resp, err := d.Client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerName)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// convert to json and write to response
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// start the container
		err = d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}
}
