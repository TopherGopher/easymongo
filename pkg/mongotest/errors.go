package mongotest

import "errors"

var (
	// ErrFailedToConnectToDockerDaemon denotes that we couldn't connect to the docker daemon
	// running locally on the machine.
	ErrFailedToConnectToDockerDaemon = errors.New("could not connect to docker daemon")
	// ErrNoAvailablePorts denotes that no ports were available for binding the docker mongo instance to
	ErrNoAvailablePorts = errors.New("no ports are available to bind the docker mongo instance to")
	// ErrMongoContainerAlreadyRunning
	ErrMongoContainerAlreadyRunning = errors.New("the mongo container is already running - an attempt was made to call it a second time")
)

type MongoTestError struct {
	err error
}

func NewMongoTestError(err error) *MongoTestError {
	return &MongoTestError{
		err: err,
	}
}
func (mte *MongoTestError) Unwrap() error {
	return mte.err
}
