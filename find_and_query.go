package easymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindAndQuery is used for Find().OneAnd() operations (e.g. Find().OneAnd().Replace())
type FindAndQuery struct {
	result interface{}
	*Query
	skip                     *int64
	limit                    *int64
	allowDiskUse             *bool
	bypassDocumentValidation *bool
	upsert                   *bool
	projection               interface{}
	arrayFilters             *options.ArrayFilters
}

// Projection limits the information returned from the query.
// Whitelist fields using:     `bson.M{"showThisField": 1}`
// Blacklist fields using:     `bson.M{"someBigStructToHide": 0}`
func (q *FindAndQuery) Projection(projectionQuery interface{}) *FindAndQuery {
	q.projection = projectionQuery
	return q
}

// Upsert specifies that if a document doesn't exist that matches the update filter,
// then a new document will be created as a result of this query run.
func (q *FindAndQuery) Upsert() *FindAndQuery {
	t := true
	q.upsert = &t
	return q
}

// ArrayFilters is used to hold filters for the array filters CRUD option. If a registry is nil, bson.DefaultRegistry
// will be used when converting the filter interfaces to BSON.
// TODO: ArrayFilters helpers
func (q *FindAndQuery) ArrayFilters(o *options.ArrayFilters) *FindAndQuery {
	q.arrayFilters = o
	return q
}

// TODO: BypassDocumentValidation options docs
func (q *FindAndQuery) BypassDocumentValidation() *FindAndQuery {
	t := true
	q.bypassDocumentValidation = &t
	return q
}

// Skip sets the skip value to bypass the given number of entries
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindAndQuery) Skip(skip int) *FindAndQuery {
	s64 := int64(skip)
	q.skip = &s64
	return q
}

// Limit sets the max value of responses to return when executing the query.
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindAndQuery) Limit(limit int) *FindAndQuery {
	// TODO: What happens with negative limits in the mongo-driver?
	if limit < 0 {
		limit = 0
	}
	l64 := int64(limit)
	q.limit = &l64
	return q
}

// Timeout uses the provided duration to set a timeout value using
// a context. The timeout clock begins upon query execution (e.g. calling .All()),
// not at time of calling Timeout().
func (q *FindAndQuery) Timeout(d time.Duration) *FindAndQuery {
	q.Query.setTimeout(d)
	return q
}

func (q *FindAndQuery) findOneAndUpdateOptions() *options.FindOneAndUpdateOptions {
	o := &options.FindOneAndUpdateOptions{
		ArrayFilters:             q.arrayFilters,
		BypassDocumentValidation: q.bypassDocumentValidation,
		Collation:                q.collation,
		MaxTime:                  q.timeout,
		Projection:               q.projection,
		Upsert:                   q.upsert,
	}
	if q.hintIndices != nil {
		o.Hint = *q.hintIndices
	}
	if q.sortFields != nil {
		o.Sort = *q.sortFields
	}
	// o.ReturnDocument can be either options.Before or options.After
	return o
}

// Update ends up running `findAndModify()` - if you do not need the result
// object, consider running `collection.UpdateOne()` instead.
func (q *FindAndQuery) Update(updateQuery interface{}) (err error) {
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := q.findOneAndUpdateOptions()
	err = mongoColl.FindOneAndUpdate(ctx, q.filter, updateQuery, opts).Decode(q.result)
	return err
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

func (q *FindAndQuery) findOneAndDeleteOptions() *options.FindOneAndDeleteOptions {
	o := &options.FindOneAndDeleteOptions{
		Collation:  q.collation,
		MaxTime:    q.timeout,
		Projection: q.projection,
	}
	if q.hintIndices != nil {
		o.Hint = *q.hintIndices
	}
	if q.sortFields != nil {
		o.Sort = *q.sortFields
	}
	return o
}

// Delete ultimately ends up running `findOneAndDelete()`. If you do not need the existing
// value/object prior to the deletion, it is recommended to instead run `Collection().Delete().One()`
func (q *FindAndQuery) Delete() (err error) {
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := q.findOneAndDeleteOptions()
	err = mongoColl.FindOneAndDelete(ctx, q.filter, opts).Decode(q.result)
	return err
}
