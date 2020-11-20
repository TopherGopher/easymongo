package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFind(t *testing.T) {
	setup(t)
	defer teardown(t)
	is := assert.New(t)
	// Create some test data
	dbName := "batman_archive"
	collName := "enemies"
	coll := conn.GetDatabase(dbName).C(collName)
	enemies := []enemy{
		0: {ID: primitive.NewObjectID(), Name: "The Joker"},
		1: {ID: primitive.NewObjectID(), Name: "Superman (depending on the day)"},
		2: {Name: "Poison Ivy"},
	}
	ids, err := coll.Insert().Many(enemies)
	is.NoError(err, "Couldn't setup the collection for the Find test")
	is.Len(ids, 3)

	t.Run("Find().One()", func(t *testing.T) {
		is := assert.New(t)
		e := enemy{}
		expectedName := "The Joker"
		err = coll.Find(bson.M{"name": expectedName}).One(&e)
		is.NoError(err, "Couldn't Find.One() the name '%s'", expectedName)
		is.Equal(expectedName, e.Name, "Returned object appears unpopulated")
	})
	t.Run("Find().Many()", func(t *testing.T) {
		is := assert.New(t)
		enemies := []enemy{}
		expectedName := "The Joker"
		err = coll.Find(bson.M{"name": expectedName}).Many(&enemies)
		is.NoError(err, "Couldn't Find.Many() the name '%s'", expectedName)
		is.Len(enemies, 1)
		is.Equal(expectedName, enemies[0].Name, "Returned object appears unpopulated")

		err = coll.Find(bson.M{}).Many(&enemies)
		is.NoError(err, "Failed to Find.Many() for all documents in collection", expectedName)
		is.Len(enemies, 3)
	})
}
