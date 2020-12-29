package easymongo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tophergopher/easymongo"
	"github.com/tophergopher/mongotest"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var conn *mongotest.TestConnection

func setup(t *testing.T) {
	is := assert.New(t)
	var err error
	conn, err = mongotest.NewTestConnection(true)
	is.NoError(err, "Could not stand up test database connection")
	if conn != nil {
		t.Cleanup(func() { teardown(t) })
	} else {
		t.FailNow()
	}
}
func teardown(t *testing.T) {
	if conn != nil {
		err := conn.KillMongoContainer()
		if err != nil {
			t.Errorf("Could not tear down test mongo container after test run: %v", err)
		}
	}
}

// enemy is a test struct for exercising unit tests
type enemy struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name"`
	Notes         string             `bson:"notes"`
	LastEncounter *time.Time         `bson:"lastEncounter"`
}

// Create some test data
func createBatmanArchive(t *testing.T) *easymongo.Collection {
	is := assert.New(t)
	dbName := "batman_archive"
	collName := "enemies"
	coll := conn.Database(dbName).C(collName)
	enemies := []enemy{
		0: {ID: primitive.NewObjectID(), Name: "The Joker", Notes: "Follow-up about his scars."},
		1: {ID: primitive.NewObjectID(), Name: "Superman", Notes: "Enemy status depends on the day - we are enemies every day on Wednesday from 4-5:30pm."},
		2: {ID: primitive.NewObjectID(), Name: "Poison Ivy"},
		3: {ID: primitive.NewObjectID(), Name: "Two-Face", Notes: "Sometimes this guy is great, othertimes, man, what a jerk"},
		4: {ID: primitive.NewObjectID(), Name: "Edward Nigma"},
	}
	ids, err := coll.Insert().Many(enemies)
	is.NoError(err, "Couldn't setup the collection for the test")
	is.Len(ids, len(enemies), "One or more items weren't inserted successfully.")
	_, err = coll.Index("name").Ensure()
	is.NoError(err, "Couldn't ensure the name index")
	return coll
}
