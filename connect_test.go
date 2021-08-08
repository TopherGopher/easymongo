package easymongo_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tophergopher/easymongo"
)

func TestConnect(t *testing.T) {
	setup(t)
	t.Cleanup(func() { teardown(t) })
	t.Run("TestPing", func(t *testing.T) {
		is := assert.New(t)
		is.NoError(conn.Ping(), "Could not ping test instance")
	})
	t.Run("ConnectWithOptions", func(t *testing.T) {
		is := assert.New(t)
		// TODO: Rest of flags
		log := logrus.New()
		tmpConn, err := easymongo.ConnectWith(conn.MongoURI()).Flags(
			easymongo.DefaultAnywhere).Debug().Logger(log).Connect()
		is.NoError(err, "Issue connecting to the test instance using options")
		is.NoError(tmpConn.Ping(), "Could not ping the test instance")
	})
	t.Run("Connect", func(t *testing.T) {
		is := assert.New(t)
		defer func() {
			if trace := recover(); trace != nil {
				t.Errorf("Connect panicked: %v", trace)
				t.Fail()
			}
		}()
		tmpConn := easymongo.Connect(conn.MongoURI())
		is.NoError(tmpConn.Ping(), "Could not ping the test instance")
	})

	t.Run("GetCurrentConnection", func(t *testing.T) {
		is := assert.New(t)
		globalConn := easymongo.GetCurrentConnection()
		is.NotNil(conn)
		is.Equal(conn.MongoURI(), globalConn.MongoURI(), "The global connection doesn't match the current connection")
	})

	t.Run("ConnectUsingMongoClient", func(t *testing.T) {
		is := assert.New(t)
		// TODO: Actually construct a client rather than using a recycled one
		tmpConn := easymongo.ConnectWith(conn.MongoURI()).FromMongoDriverClient(conn.MongoDriverClient())
		is.NoError(tmpConn.Ping(), "Could not ping the test instance")
	})
	t.Run("Connect to a DB with different flags", func(t *testing.T) {
		is := assert.New(t)
		db := conn.DatabaseByConnectionType("my_db", easymongo.DefaultPrimary)
		is.NotNil(db, "The database shouldn't be empty")
		db = conn.DatabaseByConnectionType("my_db",
			easymongo.ReadConcernLinearizable|easymongo.ReadPreferencePrimaryPreferred|easymongo.WriteConcernW3)
		is.NotNil(db, "The database shouldn't be empty")
		db = conn.DatabaseByConnectionType("my_db", easymongo.DefaultSecondary)
		is.NotNil(db, "The database shouldn't be empty")
	})
}
