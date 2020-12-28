package easymongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteQuery stores the data necessary to execute a deletion. collection.Delete() returns an initialized DeleteQuery.
type DeleteQuery struct {
	*Query
}

// Delete helps construct and execute deletion queries.
// Use with .One and .Many
func (c *Collection) Delete(filter interface{}) *DeleteQuery {
	return &DeleteQuery{
		Query: c.query(filter),
	}
}

// deleteOptions returns the native mongo driver options.DeleteOptions using
// the provided query information.
func (dq *DeleteQuery) deleteOptions() *options.DeleteOptions {
	o := &options.DeleteOptions{
		Collation: dq.collation,
	}
	if dq.hintIndices != nil {
		o.Hint = *dq.hintIndices
	}
	return o
}

// One calls out to DeleteOne() which deletes the first entry matching the
// filter query provided to Delete().
func (dq *DeleteQuery) One() (err error) {
	ctx, cancelFunc := dq.getContext()
	defer cancelFunc()
	opts := dq.deleteOptions()
	_, err = dq.collection.mongoColl.DeleteOne(ctx, dq.filter, opts)
	if err != nil {
		return err
	} else if ctx.Err() != nil {
		return ErrTimeoutOccurred
	}
	// TODO: Handle ErrNotFound
	// if res.DeletedCount == 0 { err = ErrNotFound }
	return err
}

// Many calls out to DeleteMany() which deletes all entries matching the
// filter query provided to Delete().
func (dq *DeleteQuery) Many() (numDeleted int, err error) {
	ctx, cancelFunc := dq.getContext()
	defer cancelFunc()
	opts := dq.deleteOptions()
	res, err := dq.collection.mongoColl.DeleteMany(ctx, dq.filter, opts)
	if err != nil {
		return numDeleted, err
	} else if ctx != nil && ctx.Err() != nil {
		return numDeleted, ErrTimeoutOccurred
	}
	if res != nil {
		numDeleted = int(res.DeletedCount)
	}
	// TODO: Handle ErrNotFound
	// if res.DeletedCount == 0 { err = ErrNotFound }
	return numDeleted, err
}

// DeleteByID assumes that an ID is an ObjectID and the ID is located at _id.
func (c *Collection) DeleteByID(id primitive.ObjectID) (err error) {
	return c.Delete(bson.M{"_id": id}).One()
}

func (dq *DeleteQuery) Collation(c *options.Collation) *DeleteQuery {
	dq.Query.setCollation(c)
	return dq
}

func (dq *DeleteQuery) Hint(indexKeys ...string) *DeleteQuery {
	dq.Query.setHint(indexKeys...)
	return dq
}
