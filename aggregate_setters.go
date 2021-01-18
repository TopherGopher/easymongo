package easymongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// BatchSize overrides the batch size for how documents are returned.
func (p *AggregationQuery) BatchSize(n int) *AggregationQuery {
	x := int32(n)
	p.batchSize = &x
	return p
}

// AllowDiskUse sets a flag which allows queries to page to disk space
// should they exhaust their allotted memory.
func (p *AggregationQuery) AllowDiskUse() *AggregationQuery {
	t := true
	p.allowDiskUse = &t
	return p
}

// Collation allows users to specify language-specific rules for string comparison, such as rules for lettercase and accent marks.
// https://docs.mongodb.com/manual/reference/collation/
// TODO: Create helpers and consts for Collation
func (p *AggregationQuery) Collation(c *options.Collation) *AggregationQuery {
	p.Query.setCollation(c)
	return p
}

// Comment adds a comment to the query - when the query is executed, this
// comment can help with debugging from the logs.
func (p *AggregationQuery) Comment(comment string) *AggregationQuery {
	p.Query.setComment(comment)
	return p
}

// TODO: BypassDocumentValidation options docs
func (p *AggregationQuery) BypassDocumentValidation() *AggregationQuery {
	t := true
	p.bypassDocumentValidation = &t
	return p
}

// Hint allows a user to specify index key(s) and supplies these to
// .hint() - this can result in query optimization.
// This should either be the index name as a string or the index specification
// as a document.
// The following example would instruct mongo to use a field called 'age' as
// a look-up index.
// Mongo CLI: db.users.find().hint( { age: 1 } )
// easymongo:
// err = conn.Collection(
// 	"users").Find(bson.M{}).Hint("age").One(&userObj)
// Reference: https://docs.mongodb.com/manual/reference/operator/meta/hint/
// TODO: Support '-' prepending - shoul it be -1 or 0 as the value?
func (p *AggregationQuery) Hint(indexKeys ...string) *AggregationQuery {
	p.Query.setHint(indexKeys...)
	return p
}

// Timeout uses the provided duration to set a timeout value using
// a context. The timeout clock begins upon query execution (e.g. calling .All()),
// not at time of calling Timeout().
func (p *AggregationQuery) Timeout(d time.Duration) *AggregationQuery {
	p.Query.setTimeout(d)
	return p
}

// WithContext will consume the supplied context. A note that this is completely optional
// as Timeout will auto-create a context for you if one is not supplied.
// This is typically useful when using a context across many functions.
func (p *AggregationQuery) WithContext(ctx context.Context) *AggregationQuery {
	p.Query.setContext(&ctx)
	return p
}
