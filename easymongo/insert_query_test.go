package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}
