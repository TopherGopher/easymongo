// This confusingly named file is to check to ensure the functionality
// provided by the mongotest package is working properly. No functions should be imported
// from here.
package mongotest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongoContainer(t *testing.T) {
	var err error
	var conn *TestConnection
	t.Run("A test mongo connection can be initialized with a docker container", func(t *testing.T) {
		is := assert.New(t)
		conn, err = NewTestConnection(true)
		is.NotNil(conn, "The test singleton was not initialized")
		is.NoError(err, "The test connection could not be completed")
	})
	if conn == nil {
		t.FailNow()
	}
	t.Cleanup(func() {
		_ = conn.KillMongoContainer()
	})

	t.Run("Container is running", func(t *testing.T) {
		is := assert.New(t)
		containerID := conn.MongoContainerID()
		is.NotEmpty(containerID, "No container ID was returned - a container may not have been started")
		if len(containerID) == 0 {
			t.FailNow()
		}
		args := []string{}
		top, err := conn.dockerClient.ContainerTop(
			context.Background(),
			containerID,
			args,
		)
		is.NoError(err, "Could not top the mongo container - even though it should be running")
		is.Len(top.Processes, 1, "No processes were running in the mongo container")
		is.Contains(top.Processes[0], "mongod --bind_ip_all", "The mongo daemon was not running")
	})

	t.Run("Container can be killed", func(t *testing.T) {
		is := assert.New(t)
		containerID := conn.MongoContainerID()
		err = conn.KillMongoContainer()
		is.NoError(err, "Unable to stop mongo! Call Sheriff Bart!")
		args := []string{}
		_, err = conn.dockerClient.ContainerTop(
			context.Background(),
			containerID,
			args,
		)
		is.Error(err, "The container should be gone")
		if err == nil {
			t.FailNow()
		}
		is.Contains(err.Error(), "No such container", "There shouldn't be a docker container anymore")
	})
}

func TestGetAvailablePort(t *testing.T) {
	is := assert.New(t)
	portNumber, err := GetAvailablePort()
	is.NoError(err, "Could not find an available port")
	is.Greater(portNumber, 0, "Port number should be greater than 0")
}
