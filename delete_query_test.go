package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDelete(t *testing.T) {
	var err error
	setup(t)
	coll := createBatmanArchive(t)
	t.Cleanup(func() {
		coll.Drop()
	})

	t.Run("Delete one", func(t *testing.T) {
		is := assert.New(t)
		err = coll.Delete(bson.M{"name": "The Joker"}).One()
		is.NoError(err, "Could not delete the document")
		var e enemy
		err = coll.Find(bson.M{"name": "The Joker"}).One(&e)
		is.Error(err, "Error should be set to NotFound")
		is.Empty(e.Name, "No document should have been returned")
	})
	t.Run("Delete one by ID", func(t *testing.T) {
		is := assert.New(t)
		var e enemy
		err = coll.Find(bson.M{}).One(&e)
		is.NoError(err, "Can't find a document to delete")
		is.NotEmpty(e.ID, "Couldn't find a document to delete")
		is.NotEmpty(e.Name, "Couldn't find a document to delete")

		err = coll.DeleteByID(e.ID)
		is.NoError(err, "Could not delete the document")
		e = enemy{}
		err = coll.FindByID(e.ID, &e)
		is.Error(err, "We should not be able to find this document post-deletion")
		is.Empty(e.Name, "We should not be able to find this document post-deletion")
	})
	t.Run("Delete many", func(t *testing.T) {
		is := assert.New(t)
		filter := bson.M{}
		docCount, err := coll.Find(filter).Count()
		is.NoError(err, "Could not find any documents prior to testing deletion")
		is.NotEqual(0, docCount, "Could not find any documents prior to testing deletion")

		deletedCount, err := coll.Delete(filter).Many()
		is.NoError(err, "Could not delete all documents")
		is.Equal(docCount, deletedCount, "Could not delete all the documents in the collection")

		createBatmanArchive(t)
		filter = bson.M{"notes": bson.M{"$ne": nil}}
		docCount, err = coll.Find(filter).Count()
		is.NoError(err, "Could not find any documents prior to testing deletion")
		is.NotEmpty(docCount, "Could not find any documents prior to testing deletion")

		deletedCount, err = coll.Delete(filter).Many()
		is.NoError(err, "Could not delete all documents with filter query")
		is.Equal(docCount, deletedCount, "Could not delete all the documents in the collection")
	})
}
