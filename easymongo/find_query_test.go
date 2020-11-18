package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFind(t *testing.T) {
	setup(t)
	defer teardown(t)
	var err error
	// Create some test data
	dbName, collName := createBatmanArchiveUsingMongoDriver(t)
	coll := conn.GetDatabase(dbName).C(collName)

	t.Run("Find().One()", func(t *testing.T) {
		is := assert.New(t)
		e := enemy{}
		expectedName := "The Joker"
		err = coll.Find(bson.M{"name": expectedName}).One(&e)
		is.NoError(err, "Couldn't Find.One() the name '%s'", expectedName)
		is.Equal(expectedName, e.Name, "Returned object appears unpopulated")
	})

}
