package easymongo

import "go.mongodb.org/mongo-driver/mongo/options"

// ReplaceQuery is a helper for replacement query actions and options.
type ReplaceQuery struct {
	newObj interface{}
	*Query
}

// Replace returns a ReplaceQuery. Trying running `.One()` against this.
// This is used to replace an entire object. If you are looking to update just part of a document
// (e.g. $set a field, $inc a counter up or down, etc.) you should instead use Update().One().
func (c *Collection) Replace(filter interface{}, obj interface{}) *ReplaceQuery {
	return &ReplaceQuery{
		newObj: obj,
		Query:  c.Query(filter),
	}
}

// Execute runs the ReplaceQuery. No actions are taken until this query is run.
func (rq *ReplaceQuery) Execute() error {
	// var result *mongo.UpdateResult
	opts := options.Replace()
	// TODO: ReplaceOptions
	ctx, cancelFunc := rq.getContext()
	defer cancelFunc()
	res, err := rq.collection.mongoColl.ReplaceOne(ctx, rq.filter, rq.newObj, opts)
	if err != nil {
		return err
	}
	// TODO: ReplaceQuery ErrNotFound behavior
	_ = res
	return nil
}
