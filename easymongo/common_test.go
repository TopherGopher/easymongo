package easymongo_test

import (
	"testing"
	"time"

	"github.com/tophergopher/easymongo"

	"github.com/TopherGopher/pkg/mongotest"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var conn *mongotest.TestConnection

func setup(t *testing.T) {
	is := assert.New(t)
	var err error
	conn, err = mongotest.NewTestConnection(true)
	is.NoError(err, "Could not stand up test database connection")
}
func teardown(t *testing.T) {
	if conn != nil {
		err := conn.KillMongoContainer()
		if err != nil {
			t.Errorf("Could not tear down test mongo container: %v", err)
		}
	}
}

// enemy is a test struct for exercising unit tests
type enemy struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name"`
	LastEncounter time.Time          `bson:"lastEncounter"`
}

// Create some test data
func createBatmanArchive(t *testing.T) *easymongo.Collection {
	is := assert.New(t)
	dbName := "batman_archive"
	collName := "enemies"
	coll := conn.GetDatabase(dbName).C(collName)
	enemies := []enemy{
		0: {ID: primitive.NewObjectID(), Name: "The Joker"},
		1: {ID: primitive.NewObjectID(), Name: "Superman (depending on the day)"},
		2: {Name: "Poison Ivy"},
	}
	ids, err := coll.Insert().Many(enemies)
	is.NoError(err, "Couldn't setup the collection for the test")
	is.Len(ids, 3)
	return coll
}

// func createBatmanArchiveUsingMongoDriver(t *testing.T) (dbName, collName string) {
// 	is := assert.New(t)
// 	dbName = "batman_archive"
// 	collName = "enemies"

// 	enemies := []interface{}{
// 		enemy{ID: primitive.NewObjectID(), Name: },
// 		enemy{ID: primitive.NewObjectID(), Name: },
// 	}
// 	// Perform the insert using the MongoDriverClient to avoid potential test cross-reference issues
// 	_, err := conn.MongoDriverClient().Database(dbName).Collection(collName).InsertMany(
// 		nil, enemies)
// 	is.NoError(err, "Could not setup insert test response for collection setup")
// 	return dbName, collName
// }
