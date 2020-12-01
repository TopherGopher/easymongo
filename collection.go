package easymongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

// Collection is a helper for accessing and modifying mongo collections
type Collection struct {
	database       *Database
	collectionName string
	mongoColl      *mongo.Collection
}

// Connection returns the connection associated with this collection.
func (c *Collection) Connection() *Connection {
	return c.database.connection
}

// DefaultCtx returns the appropriate context using the default timeout specified at conneciton time.
func (c *Collection) DefaultCtx() (context.Context, context.CancelFunc) {
	return c.database.DefaultCtx()
}

// GetDatabase returns the database associated with the database/collection.
func (c *Collection) GetDatabase() *Database {
	return c.database
}

// GetCollection is a shorthand for getting a collection by name using the globally
// initialized database. Consider using db.Collection() if you wish to explicitly consume
// a given connection pool.
func GetCollection(dbName, collectionName string) *Collection {
	return GetDatabase(dbName).C(collectionName)
}

// Drop drops the collection this object is referring to.
func (c *Collection) Drop() (err error) {
	ctx, cancelFunc := c.DefaultCtx()
	defer cancelFunc()
	return c.mongoColl.Drop(ctx)
}

// MongoDriverCollection returns the native mongo driver collection object
// (should you wish to interact with it directly)
func (c *Collection) MongoDriverCollection() *mongo.Collection {
	return c.mongoColl
}

// func (c *Collection) With(s *Session) *Collection {return }

// func (c *Collection) EnsureIndexKey(key ...string) error {return }
// func (c *Collection) EnsureIndex(index Index) error {return }
// func (c *Collection) Indexes() (indexes []Index, err error) {return }
// func (c *Collection) DropIndex(key ...string) error {return }
// func (c *Collection) DropIndexName(name string) error {return }

// func (c *Collection) Repair() *Iter {return }
// func (c *Collection) Aggregate(pipeline interface{}) *Aggregation { return nil }

// func (c *Collection) NewIter(session *Session, firstBatch []bson.Raw, cursorId int64, err error) *Iter {return }
// func (c *Collection) Insert(docs ...interface{}) error {return }
// func (c *Collection) Create(info *CollectionInfo) error {return }

// Count returns the count of the number of documents in the collection.
// TODO: Should we be using estimatedCount here instead?
func (c *Collection) Count() (int, error) {
	opts := options.Count()
	ctx, cancelFunc := c.database.connection.GetDefaultTimeoutCtx()
	defer cancelFunc()
	count, err := c.mongoColl.CountDocuments(ctx, nil, opts)
	return int(count), err
}

// FindByID wraps Find, ultimately executing `findOne("_id": providedID)`
// Typically, the provided id is a pointer to a *primitive.ObjectID.
func (c *Collection) FindByID(id interface{}, result interface{}) (err error) {
	return c.Find(bson.M{"_id": id}).One(result)
}

// FindByDate is a helper for filtering documents by times using the ObjectID. This is
// typically helpful when dealing with large collections as Skip and Limit will become less performant.
func (c *Collection) FindByDate(after *time.Time, before *time.Time, additionalFilters bson.M) *FindQuery {
	q := additionalFilters
	switch {
	case before != nil && after != nil:
		q["_id"] = bson.M{
			"$lte": primitive.NewObjectIDFromTimestamp(*before),
			"$gte": primitive.NewObjectIDFromTimestamp(*after),
		}
	case before != nil:
		q["_id"] = bson.M{
			"$lte": primitive.NewObjectIDFromTimestamp(*before),
		}
	case after != nil:
		q["_id"] = bson.M{
			"$gte": primitive.NewObjectIDFromTimestamp(*after),
		}
	}

	return c.Find(q)
}

// // UpdateOne updates a single record (should it exist)
// // ErrNotFound is returned if nothing matches the update criteria.
// // c.UpdateOne(bson.M{"name": "joker"})
// func (c *Collection) UpdateOne(filter interface{}, update interface{}) *UpdateQuery {
// 	return &UpdateQuery{
// 		updateQuery: update,
// 		Query: Query{
// 			filter: filter,
// 			many:   false,
// 		},
// 	}
// }

// UpdateByID wraps UpdateOne to update a single record by ID (should the record exist).
func (c *Collection) UpdateByID(id interface{}, update interface{}) (err error) {
	return c.Update(bson.M{"_id": id}, update).One()
}

// // UpdateMany updates all matching entries to the provided query.
// // If no entries were updated, ErrNotFound is returned.
// func (c *Collection) UpdateMany(filter interface{}, update interface{}) *UpdateQuery {
// 	return c.Update(filter, update).Many()
// }

// Upsert updates the first matching document using the upsert option once .One() has been called.
// If no option overrides are necessary, consider using UpsertOne or UpsertByID.
func (c *Collection) Upsert(filter interface{}, updateQuery interface{}) *UpdateQuery {
	return c.Update(filter, updateQuery).Upsert()
}

// UpsertOne updates the first matching document using the default upsert options.
// This call is equivalent to c.Upsert(filter, updateQuery).One()
func (c *Collection) UpsertOne(filter interface{}, updateQuery interface{}) error {
	return c.Upsert(filter, updateQuery).One()
}

// UpsertByID performs an upsert style update using the updateQuery against the provided _id.
func (c *Collection) UpsertByID(id interface{}, updateQuery interface{}) (err error) {
	return c.UpsertOne(bson.M{"_id": id}, updateQuery)
}

// UpsertMany performs an update style upsert using updateMany().
// Should no documents match the query, then a new document is created.
// updateQuery is typically of some sort of $set or $push bson.M.
func (c *Collection) UpsertMany(filter interface{}, updateQuery interface{}) (matchedCount, updatedCount int, err error) {
	return c.Update(filter, updateQuery).Upsert().Many()
}

func (c *Collection) DeleteOne()                {}
func (c *Collection) DeleteMany()               {}
func (c *Collection) DeleteByID(id interface{}) {}

// ReplaceByID is a friendly helper that wraps Replace(bson.M{"_id": id}, obj).Execute()
func (c *Collection) ReplaceByID(id interface{}, obj interface{}) (err error) {
	return c.Replace(bson.M{"_id": id}, obj).Execute()
}
