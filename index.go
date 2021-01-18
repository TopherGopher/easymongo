package easymongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Index represents a helper for manipulating database indices
type Index struct {
	indexNames []string
	collection *Collection
}

// Collection returns the collection object associated with this index
func (i *Index) Collection() *Collection { return i.collection }

// TODO: Index.Load() loads an index from the database into memory
func (i *Index) Load() (err error) {
	return err
}

// Ensure ensures that an index exists.
func (i *Index) Ensure() (indexName string, err error) {
	ctx, cancel := i.collection.operationCtx()
	defer cancel()
	opts := options.CreateIndexes()
	// TODO: Index.Ensure() options
	// TODO: Support compound index
	indexName, err = i.collection.mongoColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bsonx.Doc{{Key: i.indexNames[0], Value: bsonx.Int32(1)}},
	}, opts)
	err = i.collection.handleErr(err)
	return indexName, err
}

// TODO: Index.Drop() drops an index from the database
func (i *Index) Drop() (err error) {
	return err
}
