package easymongo

import "go.mongodb.org/mongo-driver/mongo/options"

// FindQuery is a helper for finding and counting documents
type FindQuery struct {
	Query
	skip  int
	limit int
}

// Skip sets the skip value to bypass the given number of entries
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindQuery) Skip(skip int) *FindQuery {
	q.skip = skip
	return q
}

// Limit sets the max value of responses to return when executing the query.
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindQuery) Limit(limit int) *FindQuery {
	q.limit = limit
	return q
}

// Many executes the specified query using find() and unmarshals
// the result into the provided interface. Ensure interface{} is either
// a slice or a pointer to a slice.
func (q *FindQuery) Many(results interface{}) error {
	// TODO: Check kind to make sure it's a slice or map
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

// FindOneOptions generates the native mongo driver FindOneOptions from the FindQuery
func (q *FindQuery) FindOneOptions() *options.FindOneOptions {
	opts := options.FindOne()
	if len(q.comment) > 0 {
		opts.SetComment(q.comment)
	}
	if len(q.hintIndices) > 0 {
		opts.SetHint(q.hintIndices)
	}
	// if q.sortFields != 0 {
	// 	opts.SetSort(q.sortFields)
	// }
	// opts.SetAllowDiskUse(b)
	// TODO: Support the rest of the find options
	return opts
}

// FindOptions generates the native mongo driver FindOptions from the FindQuery
func (q *FindQuery) FindOptions() *options.FindOptions {
	opts := options.Find()
	if q.limit > 0 {
		opts.SetLimit(int64(q.limit))
	}
	if q.skip > 0 {
		opts.SetSkip(int64(q.skip))
	}
	if len(q.comment) > 0 {
		opts.SetComment(q.comment)
	}
	if len(q.hintIndices) > 0 {
		opts.SetHint(q.hintIndices)
	}
	// if q.sortFields != 0 {
	// 	opts.SetSort(q.sortFields)
	// }
	// opts.SetAllowDiskUse(b)
	// TODO: Support the rest of the find options
	return opts
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

// CountOptions generates the native mongo driver CountOptions from the FindQuery
func (q *FindQuery) CountOptions() *options.CountOptions {
	opts := options.Count()
	if q.limit > 0 {
		opts.SetLimit(int64(q.limit))
	}
	if q.skip > 0 {
		opts.SetSkip(int64(q.skip))
	}

	// comments are not yet supported for count queries by mongo-go-driver
	// if len(q.comment) > 0 {
	// 	opts.SetComment(q.comment)
	// }
	if len(q.hintIndices) > 0 {
		opts.SetHint(q.hintIndices)
	}
	// if q.sortFields != 0 {
	// 	opts.SetSort(q.sortFields)
	// }
	// opts.SetAllowDiskUse(b)
	// TODO: Support the rest of the find options
	return opts
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

type AndFindQuery struct {
	Query
}

// AndUpdate ends up running `findAndModify()` - if you do not need the result
// object, consider running `collection.UpdateOne()` instead.
func (q *FindQuery) AndUpdate(result interface{}, updateQuery interface{}) (err error) {
	return err
}

// AndReplace ultimately ends up running `findAndReplace()`. If you do not need the existing
// value/object, it is recommended to instead run `collection.ReplaceOne()`
func (q *FindQuery) AndReplace(result interface{}, replacementObject interface{}) (err error) {
	return err
}
