package easymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find allows a user to execute a standard find() query.
// findOne(), find() and findAnd*() is run when a user calls:
//      q.One(), q.Many(). q.FindAnd().Replace(), q.FindAnd().Update() and q.FindAnd().Delete()
// TODO: Consider using bsoncore.Doc rather than interface?
func (c *Collection) Find(filter interface{}) (q *FindQuery) {
	q = &FindQuery{
		Query: c.Query(filter),
	}
	return q
}

// FindQuery is a helper for finding and counting documents
type FindQuery struct {
	*Query
	skip         *int64
	limit        *int64
	allowDiskUse *bool
	projection   interface{}
	// findOneOpts           *options.FindOneOptions
	// findManyOpts          *options.FindManyOptions
	// findOneAndReplaceOpts *options.FindOneAndReplaceOptions
	allowPartialResults *bool
	batchSize           *int32
	// cursorType          *CursorType
	max             interface{}
	maxAwaitTime    *time.Duration
	maxTime         *time.Duration
	min             interface{}
	noCursorTimeout *bool
	oplogReplay     *bool
	returnKey       *bool
	showRecordID    *bool
	snapshot        *bool
	sort            interface{}
}

// AllowPartialResults allows an operation on a sharded cluster to return partial results (in the case that a shard is inaccessible).
// An error will not be returned should only partial results be returned.
func (q *FindQuery) AllowPartialResults() *FindQuery {
	t := true
	q.allowPartialResults = &t
	return q
}

// BatchSize sets the max batch size returned by the server each time the cursor executes.
// Mostly useful when using Iter directly.
func (q *FindQuery) BatchSize(batchSize int) *FindQuery {
	i32 := int32(batchSize)
	q.batchSize = &i32
	return q
}

// TODO: FindQuery.CursorType helpers
// func (q *FindQuery) CursorType(x *CursorType) *FindQuery {
// 	return q
// }

func (q *FindQuery) Max(x interface{}) *FindQuery {
	return q
}
func (q *FindQuery) MaxAwaitTime(x *time.Duration) *FindQuery {
	return q
}
func (q *FindQuery) Min(x interface{}) *FindQuery {
	return q
}
func (q *FindQuery) NoCursorTimeout(x *bool) *FindQuery {
	return q
}
func (q *FindQuery) OplogReplay(x *bool) *FindQuery {
	return q
}

// Projection limits the information returned from the query.
// Whitelist fields using:     `bson.M{"showThisField": 1}`
// Blacklist fields using:     `bson.M{"someBigStructToHide": 0}`
func (q *FindQuery) Projection(projectionQuery interface{}) *FindQuery {
	q.projection = projectionQuery
	return q
}
func (q *FindQuery) ReturnKey(x *bool) *FindQuery {
	return q
}
func (q *FindQuery) ShowRecordID(x *bool) *FindQuery {
	return q
}

func (q *FindQuery) Snapshot(x *bool) *FindQuery {
	return q
}

// AllowDiskUse sets a flag which allows queries to page to disk space
// should they exhaust their allotted memory.
func (q *FindQuery) AllowDiskUse() *FindQuery {
	t := true
	q.allowDiskUse = &t
	return q
}

// Skip sets the skip value to bypass the given number of entries
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindQuery) Skip(skip int) *FindQuery {
	s64 := int64(skip)
	q.skip = &s64
	return q
}

// Limit sets the max value of responses to return when executing the query.
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindQuery) Limit(limit int) *FindQuery {
	l64 := int64(limit)
	q.limit = &l64
	return q
}

// Many executes the specified query using find() and unmarshals
// the result into the provided interface. Ensure interface{} is either
// a slice or a pointer to a slice.
func (q *FindQuery) Many(results interface{}) error {
	// TODO: Check kind to make sure results is a slice or map
	opts := q.FindOptions()
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	cursor, err := q.collection.mongoColl.Find(ctx, q.filter, opts)
	if err != nil {
		return err
	}
	// TODO: Inject ErrNotFound if option specified
	return cursor.All(ctx, results)
}

// One consumes the specified query and marshals the result
// into the provided interface.
func (q *FindQuery) One(result interface{}) (err error) {
	// TODO: Check kind ot make sure this is a pointer
	opts := q.FindOneOptions()
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	err = q.collection.mongoColl.FindOne(ctx, q.filter, opts).Decode(result)
	// Depending on ErrNotFound setting, consider unsetting ErrNotFound to make consistent experience
	// FindOne.Decode() is the only mongo-go-driver function that returns ErrNotFound when
	// no result was found.
	return err
}

// Iter returns an iterator that can be used to walk over and unpack the results one at a time.
func (q *FindQuery) Iter() (iter *Iter, err error) {
	// TODO: Check kind ot make sure it's a slice or map
	opts := q.FindOptions()

	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	cursor, err := q.collection.mongoColl.Find(ctx, q.filter, opts)
	if err != nil {
		return nil, err
	}
	// TODO: Inject ErrNotFound if option specified
	iter = &Iter{
		cursor: cursor,
		query:  q,
	}
	return iter, nil
}

// Count counts the number of documents in the specified query
func (q *FindQuery) Count() (int, error) {
	opts := q.CountOptions()
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	count, err := mongoColl.CountDocuments(ctx, q.filter, opts)
	return int(count), err
}

// OneAnd consumes the specified query and marshals the result
// into the provided interface once .Replace() or .Update() are called.
func (q *FindQuery) OneAnd(result interface{}) *FindAndQuery {
	// TODO: Check kind to make sure result is a pointer
	// 	// If previousDocument is not a pointer type, then we need to bail
	// if interfaceIsUnpackable(previousDocument) {
	// 	return ErrPointerRequired
	// }
	return &FindAndQuery{
		result: result,
		Query:  q.Query,
	}
}

// FindAndQuery is used for Find().OneAnd() operations (e.g. Find().OneAnd().Replace())
type FindAndQuery struct {
	result interface{}
	*Query
}

// Update ends up running `findAndModify()` - if you do not need the result
// object, consider running `collection.UpdateOne()` instead.
func (q *FindAndQuery) Update(updateQuery interface{}) (err error) {
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := options.FindOneAndUpdate()
	// TODO: FindOneAndUpdateOptions
	res := mongoColl.FindOneAndUpdate(ctx, q.filter, updateQuery, opts)

	return res.Decode(q.result)
}

// Replace ultimately ends up running `findOneAndReplace()`. If you do not need the existing
// value/object, it is recommended to instead run `collection.Replace().Execute()`
func (q *FindAndQuery) Replace(replacementObject interface{}) (err error) {
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := options.FindOneAndReplace()
	// TODO: FindOneAndReplaceOptions
	res := mongoColl.FindOneAndReplace(ctx, q.filter, replacementObject, opts)

	return res.Decode(q.result)
}

// Delete ultimately ends up running `findOneAndDelete()`. If you do not need the existing
// value/object, it is recommended to instead run `collection.Delete().Execute()`
func (q *FindAndQuery) Delete() (err error) {
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := options.FindOneAndDelete()
	// TODO: FindOneAndDeleteOptions
	res := mongoColl.FindOneAndDelete(ctx, q.filter, opts)

	return res.Decode(q.result)
}
