package easymongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReplaceQuery is a helper for replacement query actions and options.
type ReplaceQuery struct {
	newObj interface{}
	Query
}

// FindOneAndReplaceOptions returns the options associated with this query
func (rq *ReplaceQuery) FindOneAndReplaceOptions() *options.FindOneAndReplaceOptions {
	opts := options.FindOneAndReplace()
	// TODO: Set options
	return opts
}

// Do runs the ReplaceQuery. No actions are taken until this query is run.
func (rq *ReplaceQuery) Do(previousDocument interface{}) (err error) {
	// var result *mongo.UpdateResult
	mongoColl := rq.collection.mongoColl
	ctx, cancelFunc := rq.getContext()
	defer cancelFunc()

	opts := rq.FindOneAndReplaceOptions()
	res := mongoColl.FindOneAndReplace(ctx, rq.filter, rq.newObj, opts)
	if previousDocument == nil {
		// Nothing more to be done - not unpacking
		return nil
		// TODO: If this is an uninitialized pointer (which is nil) - ensure that we continue on
	}
	// If previousDocument is not a pointer type, then we need to bail
	if interfaceIsUnpackable(previousDocument) {
		return ErrPointerRequired
	}
	return res.Decode(previousDocument)
}
