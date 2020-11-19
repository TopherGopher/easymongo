// Package mongotest provides helpers for running regressions using mongo.
// You can find helpers for:
// - ?running a database using docker
// - importing data from files
// - cleaning up a database
package mongotest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"runtime/debug"
	"strconv"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/TopherGopher/easymongo"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

// TestConnection contains helpers for creating your own tests with mongo.
type TestConnection struct {
	*easymongo.Connection
	dockerClient     *docker.Client
	logger           *logrus.Entry
	mongoContainerID string
}

// NewTestConnection is the standard method for initializing a TestConnection - it has a side-effect
// of spawning a new docker container if spinupDockerContainer is set to true.
// Note that the first time this is called on a new system, the mongo docker
// container will be pulled. Any subsequent calls on the system should succeed without
// calls to pull.
// If spinupDockerContainer is False, then no docker shenanigans occur, instead
// an attempt is made to connect to a locally running mongo instance
// (e.g. mongodb://127.0.0.1:27017).
func NewTestConnection(spinupDockerContainer bool) (*TestConnection, error) {
	// TODO: How should we be handling logging? What do other libraries typically do?
	logger := logrus.New().WithField("src", "mongotest.TestConnection")
	mongoURI := "mongodb://127.0.0.1"
	testConn := &TestConnection{
		logger: logger,
	}
	defer func() {
		if err := recover(); err != nil {
			logger.WithFields(logrus.Fields{
				"err":   err,
				"stack": string(debug.Stack()),
			}).Error("A panic occurred when trying to initialize a TestConnection")
			// Initialization crashed - ensure the mongo container is destroyed
			_ = testConn.KillMongoContainer()
		}
	}()
	if spinupDockerContainer {
		dockerClient, err := docker.NewEnvClient()
		if err != nil {
			testConn.logger.WithField("err", err).Error("Could not connect to docker daemon")
			return testConn, ErrFailedToConnectToDockerDaemon
		}
		testConn.dockerClient = dockerClient
		portNumber, err := GetAvailablePort()
		if err != nil {
			testConn.logger.WithField("err", err).Error("No ports were available to bind the test docker mongo container to")
			return testConn, ErrNoAvailablePorts
		}
		// TODO: Consider using different error types for these returns
		containerID, err := testConn.StartMongoContainer(portNumber)
		if err != nil {
			logger.WithField("err", err).Error("Could not spawn the to mongo container")
			return testConn, err
		}
		testConn.mongoContainerID = containerID
		mongoURI = fmt.Sprintf("%s:%d", mongoURI, portNumber)
	}
	conn, err := easymongo.ConnectWithOptions(mongoURI, nil)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"mongoURI": mongoURI,
		}).Error("Could not connect to mongo instance")
		return testConn, err
	}
	// Allow up to 1 second for the mongo container to come up across 5 retrie=
	numChecks := 5
	sleepTime := time.Millisecond * 200
	for i := 0; i < numChecks; i++ {
		if err = conn.Ping(); err == nil {
			// If we were able to ping the instance, we can break
			break
		}
		logger.WithFields(logrus.Fields{
			"currentRetry":      i + 1,
			"maxRetries":        numChecks,
			"sleepMilliseconds": sleepTime.Milliseconds(),
		}).Debug("Could not connect to test database - sleeping and retrying.")
		// otherwise, we need to wait a bit before checking again
		time.Sleep(sleepTime)
	}
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"mongoURI": mongoURI,
		}).Errorf("Could not ping the test mongo instance after %d checks", numChecks)
		// Try to teardown the mongo container (it might not have started)
		_ = testConn.KillMongoContainer()
		return testConn, err
	}

	testConn.Connection = conn
	return testConn, nil
}

// MongoContainerID returns the ID of the running docker container
// If no container is running, an empty string will be returned.
func (tc *TestConnection) MongoContainerID() string {
	return tc.mongoContainerID
}

// func (tc *TestConnection) ImportFromFile(filepath string) {
// 	// Open the file

// }

// GetAvailablePort returns an available port on the system.
func GetAvailablePort() (port int, err error) {
	// Create a new server without specifying a port
	// which will result in an open port being chosen
	server, err := net.Listen("tcp", "127.0.0.1:0")
	// If there's an error it likely means no ports
	// are available or something else prevented finding
	// an open port
	if err != nil {
		return 0, ErrNoAvailablePorts
	}
	defer server.Close()
	// Get the host string in the format "127.0.0.1:4444"
	hostString := server.Addr().String()
	// Split the host from the port
	_, portString, err := net.SplitHostPort(hostString)
	if err != nil {
		return 0, err
	}

	// Return the port as an int
	// TODO: This is used as a string elsewhere - consider string
	return strconv.Atoi(portString)
}

// pullMongoContainer fetches the mongo container from dockerhub
func (tc *TestConnection) pullMongoContainer(mongoImageName string) (err error) {
	// TODO: Is this better to do as an error handler?
	// Pull the initial container
	tc.logger.Info("Starting mongo docker image pull")
	rc, err := tc.dockerClient.ImagePull(context.Background(), mongoImageName, types.ImagePullOptions{})
	defer rc.Close()
	if err != nil {
		return fmt.Errorf("could not pull mongo container: %v", err)
	}
	if _, err := ioutil.ReadAll(rc); err != nil {
		return fmt.Errorf("could not pull mongo container: %v", err)
	}
	tc.logger.Info("Done pulling mongo docker image")
	return nil
}

// StartMongoContainer starts a mongo docker container
// A note that the docker daemon on the system is expected to be running
// TODO: Is there a way to spawn the docker daemon myself?
func (tc *TestConnection) StartMongoContainer(portNumber int) (containerID string, err error) {
	if len(tc.mongoContainerID) != 0 {
		return "", ErrMongoContainerAlreadyRunning
	}
	pName := fmt.Sprintf("%d/tcp", portNumber)
	containerName := fmt.Sprintf("mongo-%d", portNumber)

	mongoImageName := "registry.hub.docker.com/library/mongo:latest"

	containerResp, err := tc.dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			// TODO: Allow setting mongo version
			Image: mongoImageName,
			Labels: map[string]string{
				"mongotest": "regression",
			},
			Tty:       true,
			OpenStdin: true,
			ExposedPorts: nat.PortSet{
				nat.Port(pName): {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port("27017/tcp"): []nat.PortBinding{
					{
						HostIP:   "127.0.0.1",
						HostPort: pName,
					},
				},
			},
		},
		// TODO: Does this config also need to be specified?
		&network.NetworkingConfig{},
		containerName)
	if err != nil && docker.IsErrNotFound(err) {
		// The image didn't exist locally - go grab it
		if err = tc.pullMongoContainer(mongoImageName); err != nil {
			// The pull didn't succeed, bail
			tc.logger.WithField("err", err).Error("Could not pull the docker container")
			return "", err
		}
		// Now that the pull is complete, we can try to call start again
		return tc.StartMongoContainer(portNumber)
	} else if err != nil {
		tc.logger.WithField("err", err).Error("Could not create the docker container")
		return "", err
	}
	containerID = containerResp.ID
	tc.mongoContainerID = containerID

	err = tc.dockerClient.ContainerStart(
		context.Background(),
		containerID,
		types.ContainerStartOptions{})
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"containerID": containerID,
			"err":         err,
		}).Error("Could not start the docker container")
		return containerID, err
	}
	tc.logger.WithFields(
		logrus.Fields{
			"containerName":      containerName,
			"containerMongoPort": portNumber,
			"containerID":        containerID,
		},
	).Info("Successfully spawned mongo docker test container.")
	return containerID, err
}

// KillMongoContainer tears down the specified container
// TODO: Is there some nifty hook I could use that allows me
// to always call this as the scope of a test exits?
func (tc *TestConnection) KillMongoContainer() (err error) {
	if len(tc.mongoContainerID) == 0 {
		// No container was ever launched, nothing to be done
		return nil
	}
	err = tc.dockerClient.ContainerRemove(context.Background(),
		tc.mongoContainerID,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"err":         err,
			"containerID": tc.mongoContainerID,
		}).Error("Could not remove container")
		return err
	}
	tc.logger.WithField("containerID", tc.mongoContainerID).Debug(
		"Successfully removed container")
	// Once removed - unset the container ID
	tc.mongoContainerID = ""
	return nil
}

// TODO: DropAllDatabases
