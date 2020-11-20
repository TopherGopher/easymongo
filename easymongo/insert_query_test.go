package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsert(t *testing.T) {
	setup(t)
	t.Cleanup(func() { teardown(t) })

	dbName := "olympus"
	collName := "greek_gods"
	type greekGod struct {
		Name string `bson:"name"`
	}
	coll := conn.GetDatabase(dbName).C(collName)
	t.Run("Insert().One()", func(t *testing.T) {
		is := assert.New(t)
		id, err := coll.Insert().One(greekGod{Name: "Zeus"})
		is.NoError(err, "Could not insert a single document")
		is.NotNil(id)
		// Ensure that it actually is stored in the DB
		var zeusLookup greekGod
		err = coll.FindByID(id, &zeusLookup)
		is.NoError(err, "Could not look-up the zeus record in the DB")
		is.Equal("Zeus", zeusLookup.Name, "The record appears to have improper information in it")
	})
	// TODO: Test IsDup failure path
	t.Run("Insert().Many()", func(t *testing.T) {
		is := assert.New(t)
		ids, err := coll.Insert().Many(&[]greekGod{{Name: "Hera"}, {Name: "Hades"}})
		is.NoError(err, "Insert().Many() failed")
		is.Len(ids, 2, "2 IDs should have been returned from the insert")
		var godLookup []greekGod
		err = coll.Find(bson.M{}).Many(&godLookup)
		is.NoError(err, "Could not find any documents")
		is.Len(godLookup, 3, "After the previous test and this test, there should be 3 documents")
	})
}
