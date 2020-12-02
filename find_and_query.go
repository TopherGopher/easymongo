package easymongo

import "go.mongodb.org/mongo-driver/mongo/options"

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

func (q *FindAndQuery) findOneAndUpdateOptions() *options.FindOneAndUpdateOptions {
	return &options.FindOneAndUpdateOptions{
		ArrayFilters:             q.arrayFilters,
		BypassDocumentValidation: q.bypassDocumentValidation,
		Collation:                q.collation,
		Hint:                     q.hintIndices,
		MaxTime:                  q.timeout,
		Projection:               q.projection,
		Sort:                     q.sortFields,
		Upsert:                   q.upsert,
	}
}

// Update ends up running `findAndModify()` - if you do not need the result
// object, consider running `collection.UpdateOne()` instead.
func (q *FindAndQuery) Update(updateQuery interface{}) (err error) {
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	opts := q.findOneAndUpdateOptions()
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
