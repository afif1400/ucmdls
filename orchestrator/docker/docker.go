package docker

import "github.com/docker/docker/client"

// Docker is the interface to the Docker API.
type Docker struct {
	// Client is the Docker API client.
	Client *client.Client
}

func NewDockerClient() (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Docker{Client: cli}, nil
}
