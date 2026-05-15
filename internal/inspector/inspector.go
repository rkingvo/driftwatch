package inspector

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ContainerInfo holds the runtime state of a container relevant to drift detection.
type ContainerInfo struct {
	Name   string
	Image  string
	Env    map[string]string
	Labels map[string]string
	Status string
}

// Client wraps the Docker client for container inspection.
type Client struct {
	docker *client.Client
}

// New creates a new inspector Client using the default Docker environment.
func New() (*Client, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("inspector: failed to create docker client: %w", err)
	}
	return &Client{docker: dc}, nil
}

// Inspect retrieves runtime info for the container with the given name or ID.
func (c *Client) Inspect(ctx context.Context, nameOrID string) (*ContainerInfo, error) {
	raw, err := c.docker.ContainerInspect(ctx, nameOrID)
	if err != nil {
		return nil, fmt.Errorf("inspector: inspect %q: %w", nameOrID, err)
	}
	return containerInfoFromRaw(raw), nil
}

// InspectAll returns runtime info for all running containers.
func (c *Client) InspectAll(ctx context.Context) ([]*ContainerInfo, error) {
	list, err := c.docker.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("inspector: list containers: %w", err)
	}

	var infos []*ContainerInfo
	for _, c2 := range list {
		raw, err := c.docker.ContainerInspect(ctx, c2.ID)
		if err != nil {
			return nil, fmt.Errorf("inspector: inspect %q: %w", c2.ID, err)
		}
		infos = append(infos, containerInfoFromRaw(raw))
	}
	return infos, nil
}

// Close releases resources held by the Docker client.
func (c *Client) Close() error {
	return c.docker.Close()
}

func containerInfoFromRaw(raw types.ContainerJSON) *ContainerInfo {
	name := raw.Name
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	env := make(map[string]string)
	for _, e := range raw.Config.Env {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				env[e[:i]] = e[i+1:]
				break
			}
		}
	}

	return &ContainerInfo{
		Name:   name,
		Image:  raw.Config.Image,
		Env:    env,
		Labels: raw.Config.Labels,
		Status: raw.State.Status,
	}
}
