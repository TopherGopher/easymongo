package easymongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AggregationQuery struct {
	*Query
	allowDiskUse             *bool
	batchSize                *int32
	bypassDocumentValidation *bool
	// Collation:                nil,
	// Comment:                  nil,
	// Hint
}

// Aggregate begins an aggregationQuery pipeline.
// The pipeline will be executed on a call to coll.Aggregate().All() or coll.Agregate().One()
func (c *Collection) Aggregate(pipeline interface{}) *AggregationQuery {
	return &AggregationQuery{
		Query: c.query(pipeline),
	}
}

// Cursor executes the query using the provided options and returns the
// mongo.Cursor that can be worked with directly. This is typically useful
// when returning large numbers of results. If you just need to get at the documents without
// iterating, call .One() or .All()
func (p *AggregationQuery) Cursor() (*mongo.Cursor, error) {
	coll := p.collection.mongoColl
	ctx, cancelFunc := p.getContext()
	defer cancelFunc()
	opts := p.aggregateOptions()
	cursor, err := coll.Aggregate(ctx, p.filter, opts)
	return cursor, p.collection.handleErr(err)
}

func (p *AggregationQuery) aggregateOptions() *options.AggregateOptions {
	opts := &options.AggregateOptions{
		AllowDiskUse:             p.allowDiskUse,
		BatchSize:                p.batchSize,
		BypassDocumentValidation: p.bypassDocumentValidation,
		Collation:                p.Query.collation,
		MaxTime:                  p.Query.timeout,
		Comment:                  p.Query.comment,
	}
	if p.hintIndices != nil {
		opts.Hint = *p.Query.hintIndices
	}
	return opts
}

// // FetchRawResult is a somewhat odd but handy helper. Sometimes it's tricky to understand why
// // values aren't unpacking as anticipated into your struct. In this case, call this function
// // to see the key/value pairs that are being returned. I haven't found a good way to display
// // the raw original result returned from the query, but it's sometimes helpful to view these key/value pairs.
// TODO: Figure out how to get at the raw result that comes back from mongo.
// func (p *AggregationQuery) FetchRawResult() (string, error) {
// 	arr := RawMongoResult{}
// 	err := p.collection.Aggregate(p.filter).All(&arr)
// 	if err != nil {
// 		return "", err
// 	}
// 	return arr.String(), err
// }

// All executes the aggregation and returns the resultant output to the provided result object.
func (p *AggregationQuery) All(result interface{}) error {
	cursor, err := p.Cursor()
	if err != nil {
		return err
	}
	ctx, cancelFunc := p.getContext()
	defer cancelFunc()
	if err = cursor.All(ctx, result); err != nil {
		return err
	}

	return p.collection.handleErr(err)
}

// One executes the aggregation and returns the first result to the provided result object.
func (p *AggregationQuery) One(result interface{}) error {
	cursor, err := p.Cursor()
	if err != nil {
		return err
	}
	ctx, cancelFunc := p.getContext()
	defer cancelFunc()
	if found := cursor.Next(ctx); !found {
		// Move the cursor
		err = ErrNoDocuments
		return err
	}
	if err = cursor.Decode(result); err != nil {
		// Decode the result
		return err
	}
	return p.collection.handleErr(err)
}
