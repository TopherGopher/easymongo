package easymongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateQuery helps construct and execute update queries
type UpdateQuery struct {
	updateQuery              interface{}
	upsert                   *bool
	bypassDocumentValidation *bool
	arrayFilters             *options.ArrayFilters
	Query
}

// Update
// todo: update docs
func (c *Collection) Update(filter interface{}, update interface{}) *UpdateQuery {
	return &UpdateQuery{
		updateQuery: update,
		Query: Query{
			filter: filter,
		},
	}
}

// ArrayFilters is used to hold filters for the array filters CRUD option. If a registry is nil, bson.DefaultRegistry
// will be used when converting the filter interfaces to BSON.
// TODO: ArrayFilters helpers
func (uq *UpdateQuery) ArrayFilters(o *options.ArrayFilters) *UpdateQuery {
	uq.arrayFilters = o
	return uq
}

// Upsert specifies that if a document doesn't exist that matches the update filter,
// then a new document will be created as a result of this query run.
func (uq *UpdateQuery) Upsert() *UpdateQuery {
	t := true
	uq.upsert = &t
	return uq
}

// TODO: BypassDocumentValidation options docs
func (uq *UpdateQuery) BypassDocumentValidation() *UpdateQuery {
	t := true
	uq.bypassDocumentValidation = &t
	return uq
}

// UpdateOptions returns the native mongo driver options.UpdateOptions using
// the provided query information.
func (uq *UpdateQuery) UpdateOptions() *options.UpdateOptions {
	return &options.UpdateOptions{
		ArrayFilters:             uq.arrayFilters,
		BypassDocumentValidation: uq.bypassDocumentValidation,
		Collation:                uq.collation,
		Hint:                     uq.hintIndices,
		Upsert:                   uq.upsert,
	}
}

// One runs the UpdateQuery against the first matching document.
// No actions are taken until this function is called.
func (uq *UpdateQuery) One() (err error) {
	var result *mongo.UpdateResult
	mongoColl := uq.collection.mongoColl
	ctx, cancelFunc := uq.getContext()
	defer cancelFunc()
	opts := uq.UpdateOptions()
	result, err = mongoColl.UpdateOne(ctx, uq.filter, uq.updateQuery, opts)
	if err == nil && result.MatchedCount == 0 {
		// TODO: Inject ErrNotFound
	}
	// matchedCount = int(result.MatchedCount)
	// updatedCount = int(result.ModifiedCount)

	return err
}

// Many runs the UpdateQuery against all matching documents.
// No actions are taken until this function is called.
func (uq *UpdateQuery) Many() (matchedCount, updatedCount int, err error) {
	var result *mongo.UpdateResult
	mongoColl := uq.collection.mongoColl
	ctx, cancelFunc := uq.getContext()
	defer cancelFunc()
	opts := uq.UpdateOptions()
	result, err = mongoColl.UpdateMany(ctx, uq.filter, uq.updateQuery, opts)
	if err == nil && result.MatchedCount == 0 {
		// TODO: Inject ErrNotFound
	}
	matchedCount = int(result.MatchedCount)
	updatedCount = int(result.ModifiedCount)

	return matchedCount, updatedCount, err
}
