package easymongo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdate(t *testing.T) {
	setup(t)
	coll := createBatmanArchive(t)
	t.Cleanup(func() {
		coll.Drop()
	})

	t.Run("Update one", func(t *testing.T) {
		is := assert.New(t)
		var e enemy
		now := time.Now()
		err := coll.Update(bson.M{"name": "The Joker"}, bson.M{
			"$set": bson.M{"lastEncounter": now}}).One()
		is.NoError(err, "Could not update the timestamp for a single document")
		err = coll.Find(bson.M{"name": "The Joker"}).One(&e)
		is.NoError(err, "Could not find the document after updating")
		is.Equal(now.Unix(), e.LastEncounter.Unix(), "The document does not appear to be updated")
	})
	t.Run("Update one by ID", func(t *testing.T) {
		is := assert.New(t)
		var e enemy
		notes := "The darkest night..."
		err := coll.Find(bson.M{"_id": bson.M{"$ne": nil}}).One(&e)
		is.NoError(err, "Unable to find an object to update by ID")
		err = coll.UpdateByID(e.ID, bson.M{"$set": bson.M{"notes": notes}})
		is.NoError(err, "Could not update the object by ID")
		err = coll.Find(bson.M{"_id": bson.M{"$ne": nil}}).One(&e)
		is.NoError(err, "Unable to find an object to update by ID")
		is.Equal(notes, e.Notes, "The document does not appear to be updated")
	})
	t.Run("Update many", func(t *testing.T) {
		is := assert.New(t)
		now := time.Now()
		var enemies []enemy
		filter := bson.M{"name": nil}
		update := bson.M{"$set": bson.M{"lastEncounter": now}}
		matchedCount, updatedCount, err := coll.Update(filter, update).Upsert().All()
		is.NoError(err, "Could not update many")
		is.Equal(0, matchedCount)
		is.Equal(1, updatedCount)

		err = coll.Find(filter).All(&enemies)
		is.NoError(err, "Unable to find an object to update by ID")
		is.Len(enemies, updatedCount)
	})
}
