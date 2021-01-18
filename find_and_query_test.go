package easymongo_test

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFindAnd(t *testing.T) {
	setup(t)
	coll := createBatmanArchive(t)
	t.Run("FindAnd().Update()", func(t *testing.T) {
		is := assert.New(t)
		var enemyBefore, enemyAfter enemy
		filter := bson.M{
			"name": "Edward Nigma",
		}
		notes := "This man is a real mystery."
		set := bson.M{
			"$set": bson.M{"notes": notes},
		}
		err := coll.Find(filter).OneAnd(&enemyBefore).Update(set)
		is.NoError(err, "Could not find a document then update it")
		is.False(enemyBefore.ID.IsZero(), "The ID should not be empty from the lookup")
		is.Equal("Edward Nigma", enemyBefore.Name, "The document doesn't appear to be properly populated")

		err = coll.Find(filter).One(&enemyAfter)
		is.NoError(err, "Couldn't look up document after FindOneAndUpdate")
		is.Equal(notes, enemyAfter.Notes, "The update to the document doesn't appear to be working")
	})
	t.Run("FindAnd().Replace()", func(t *testing.T) {
		is := assert.New(t)
		var enemyBefore, enemyAfter, enemyLookup enemy
		filter := bson.M{
			"name": "Edward Nigma",
		}

		err := coll.Find(filter).One(&enemyLookup)
		is.NoError(err, "Couldn't look up document prior to FindOneAndReplace")

		now := time.Now()
		enemyAfter = enemy{
			ID:            enemyLookup.ID,
			Name:          "The Riddler",
			LastEncounter: &now,
		}
		// The name changed - we can query by ID now
		filter = bson.M{
			"_id": enemyLookup.ID,
		}
		err = coll.Find(filter).OneAnd(&enemyBefore).Replace(enemyAfter)
		is.NoError(err, "Could not find a document then update it")
		is.False(enemyBefore.ID.IsZero(), "The ID should not be empty from the lookup")
		is.Equal("Edward Nigma", enemyBefore.Name, "The before document doesn't appear to be properly populated")

		err = coll.Find(filter).One(&enemyLookup)
		is.NoError(err, "Couldn't look up document after FindOneAndReplace")
		is.Equal(now.Unix(), enemyLookup.LastEncounter.Unix(), "The timestamp on the document didn't update")
		is.Equal("The Riddler", enemyLookup.Name, "The replace to the document doesn't appear to have worked")
	})
	t.Run("FindAnd().Delete()", func(t *testing.T) {
		is := assert.New(t)
		var enemyBefore, enemyAfter enemy
		filter := bson.M{
			"name": "The Riddler",
		}
		err := coll.Find(filter).OneAnd(&enemyBefore).Delete()
		is.NoError(err, "Could not find a document then delete it")
		is.False(enemyBefore.ID.IsZero(), "The ID should not be empty from the lookup before")
		is.Equal("The Riddler", enemyBefore.Name, "The document doesn't appear to be properly populated")

		err = coll.Find(filter).One(&enemyAfter)
		is.Equal(mongo.ErrNoDocuments, err, "The document shouldn't exist after the deletion")
		is.True(enemyAfter.ID.IsZero(), "The ID should be empty after deletion")
		is.Nil(enemyAfter.LastEncounter, "The update to the document doesn't appear to be working")
	})
}
