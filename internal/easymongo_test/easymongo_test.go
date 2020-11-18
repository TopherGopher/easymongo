// Package easymongo_test contains helpers for validating that easymongo is working correctly.
// Not to be confused with mongotest, which contains the tools you're probably looking for.
package easymongo_test

import (
	"context"
	"testing"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/TopherGopher/easymongo"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
)

func setupConnection(t *testing.T) {
	is := assert.New(t)
	dockerClient, err := docker.NewEnvClient()
	is.NoError(err, "Ensure that your docker daemon is running")
	containerName := "mongo-regression"
	containerResp, err := dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "mongo:4.0.17",
			ExposedPorts: nat.PortSet{
				"27017/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"27017/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "27017",
					},
				},
			},
		},
		&network.NetworkingConfig{},
		containerName)
	is.NoError(err, "Could not create the docker container")

	err = dockerClient.ContainerStart(
		context.Background(),
		containerResp.ID,
		types.ContainerStartOptions{})
	is.NoError(err)

	// Now connect to it
	conn, err := easymongo.Connect("mongodb://127.0.0.1:27017")
	is.NoError(err, "Could not connect to mongo instance")
	is.NoError(conn.Ping(), "Could not ping the mongo instance")
}

// func cleanupDatabase() {}
func TestConnect(t *testing.T) {
	t.Run("Test ConnectWithOptions", func(t *testing.T) { t.SkipNow() })
	t.Run("Valid Global Connection", func(t *testing.T) { t.SkipNow() })
	t.Run("Test Connect", func(t *testing.T) { t.SkipNow() })
	t.Run("Test ConnectUsingMongoClient", func(t *testing.T) { t.SkipNow() })
}

func TestInsert(t *testing.T) {}
