package easymongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateQuery helps construct and execute update queries
type UpdateQuery struct {
	updateQuery interface{}
	upsert      bool
	Query
}

// UpdateOptions returns the native mongo driver options.UpdateOptions using
// the provided query information.
func (uq *UpdateQuery) UpdateOptions() *options.UpdateOptions {
	opts := options.Update()
	// TODO: Support these other options
	// opts.SetArrayFilters(af)
	// opts.SetBypassDocumentValidation(b)
	// opts.SetCollation(c)
	if len(uq.hintIndices) > 0 {
		// TODO: Is hintIndices in the correct format?
		opts.SetHint(uq.hintIndices)
	}
	if uq.upsert {
		opts.SetUpsert(true)
	}
	return opts
}

// Do runs the UpdateQuery. No actions are taken until this query is run.
func (uq *UpdateQuery) Do() (matchedCount, updatedCount int, err error) {
	var result *mongo.UpdateResult
	mongoColl := uq.collection.mongoColl
	ctx, cancelFunc := uq.getContext()
	defer cancelFunc()
	opts := uq.UpdateOptions()
	if uq.many {
		result, err = mongoColl.UpdateMany(ctx, uq.filter, uq.updateQuery, opts)
		if err == nil && result.MatchedCount == 0 {
			// TODO: Inject ErrNotFound
		}
	} else {
		result, err = mongoColl.UpdateOne(ctx, uq.filter, uq.updateQuery, opts)
		if err == nil && result.MatchedCount == 0 {
			// TODO: Inject ErrNotFound
		}
	}
	matchedCount = int(result.MatchedCount)
	updatedCount = int(result.ModifiedCount)

	return matchedCount, updatedCount, err
}
