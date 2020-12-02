package easymongo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFind(t *testing.T) {
	setup(t)

	coll := createBatmanArchive(t)

	t.Run("Find().One()", func(t *testing.T) {
		is := assert.New(t)
		e := enemy{}
		expectedName := "The Joker"
		err := coll.Find(bson.M{"name": expectedName}).One(&e)
		is.NoError(err, "Couldn't Find.One() the name '%s'", expectedName)
		is.Equal(expectedName, e.Name, "Returned object appears unpopulated")
	})
	t.Run("Find().Many()", func(t *testing.T) {
		is := assert.New(t)
		enemies := []enemy{}
		expectedName := "The Joker"
		err := coll.Find(bson.M{"name": expectedName}).Comment(
			"Isn't this a fun query?").BatchSize(5).Projection(
			bson.M{"name": 1}).Hint("name").Sort(
			"-name").Skip(0).Limit(0).Timeout(time.Hour).Many(&enemies)
		is.NoError(err, "Couldn't Find.Many() the name '%s'", expectedName)
		is.Len(enemies, 1)
		if len(enemies) < 1 {
			t.FailNow()
		}
		is.Equal(expectedName, enemies[0].Name, "Returned object appears unpopulated")

		err = coll.Find(bson.M{}).Many(&enemies)
		is.NoError(err, "Failed to Find.Many() for all documents in collection", expectedName)
		is.Len(enemies, 3)
	})
}
