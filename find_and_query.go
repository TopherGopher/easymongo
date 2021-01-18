package easymongo

import (
	"fmt"
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
	returnDocument           options.ReturnDocument
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
		// Pass through the relevant options from Find
		skip:         q.skip,
		limit:        q.limit,
		allowDiskUse: q.allowDiskUse,
		projection:   q.projection,
		// Default return to options.Before
		returnDocument: options.Before,
		// allowPartialResults: q.allowPartialResults,
		// batchSize:           q.batchSize,
		// maxTime:             q.maxTime,
		// cursorType:          q.cursorType,
		// max:                 q.max,
		// maxAwaitTime:        q.maxAwaitTime,
		// min:                 q.min,
		// noCursorTimeout:     q.noCursorTimeout,
		// oplogReplay:         q.oplogReplay,
		// returnKey:           q.returnKey,
		// showRecordID: q.showRecordID,
		// snapshot:     q.snapshot,
		// sort:         q.sort,
	}
}

// Projection limits the information returned from the query.
// Whitelist fields using:     `bson.M{"showThisField": 1}`
// Blacklist fields using:     `bson.M{"someBigStructToHide": 0}`
func (q *FindAndQuery) Projection(projectionQuery interface{}) *FindAndQuery {
	q.projection = projectionQuery
	return q
}

// Upsert sets an option to specify that if a document doesn't exist which matches the update filter,
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

// ReturnDocumentAfterModification specifies the object should be returned after modification is complete.
// By default, the document is returned before the query.
func (q *FindAndQuery) ReturnDocumentAfterModification() *FindAndQuery {
	q.returnDocument = options.After
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
		ReturnDocument:           &q.returnDocument,
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

func (q *FindAndQuery) findOneAndReplaceOptions() *options.FindOneAndReplaceOptions {
	o := &options.FindOneAndReplaceOptions{
		BypassDocumentValidation: q.bypassDocumentValidation,
		Collation:                q.collation,
		MaxTime:                  q.timeout,
		Projection:               q.projection,
		Upsert:                   q.upsert,
		ReturnDocument:           &q.returnDocument,
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

// Update ends up running `findAndModify()` to update (and return) the first matching document
// If you do not need the result object, consider running `collection.UpdateOne()` instead.
// mongo.ErrNoDocuments is returned in the case that nothing matches the specified query.
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
// mongo.ErrNoDocuments is returned in the case that nothing matches the specified query.
func (q *FindAndQuery) Replace(replacementObject interface{}) (err error) {
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := q.findOneAndReplaceOptions()
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
// mongo.ErrNoDocuments is returned in the case that nothing matches the specified query.
func (q *FindAndQuery) Delete() (err error) {
	if q.returnDocument == options.After {
		// Explicitly fail if the user is attempting to get the document after it was deleted
		return fmt.Errorf("options.After is not compatible with FindAnd().Delete() as the document will always not be found after deletion")
	}
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := q.findOneAndDeleteOptions()
	err = mongoColl.FindOneAndDelete(ctx, q.filter, opts).Decode(q.result)
	return err
}
