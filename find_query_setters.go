package easymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllowDiskUse sets a flag which allows queries to page to disk space
// should they exhaust their allotted memory.
func (q *FindQuery) AllowDiskUse() *FindQuery {
	t := true
	q.allowDiskUse = &t
	return q
}

// AllowPartialResults allows an operation on a sharded cluster to return partial results (in the case that a shard is inaccessible).
// An error will not be returned should only partial results be returned.
func (q *FindQuery) AllowPartialResults() *FindQuery {
	t := true
	q.allowPartialResults = &t
	return q
}

// Collation allows users to specify language-specific rules for string comparison, such as rules for lettercase and accent marks.
// https://docs.mongodb.com/manual/reference/collation/
// TODO: Create helpers and consts for Collation
func (q *FindQuery) Collation(c *options.Collation) *FindQuery {
	q.Query.Collation(c)
	return q
}

// Comment adds a comment to the query - when the query is executed, this
// comment can help with debugging from the logs.
func (q *FindQuery) Comment(comment string) *FindQuery {
	q.Query.Comment(comment)
	return q
}

// BatchSize sets the max batch size returned by the server each time the cursor executes.
// Mostly useful when using Iter directly.
func (q *FindQuery) BatchSize(batchSize int) *FindQuery {
	i32 := int32(batchSize)
	q.batchSize = &i32
	return q
}

// Sort accepts a list of strings to use as sort fields.
// Prepending a field name with a '-' denotes descending sorting
// e.g. "-name" would sort the "name" field in descending order
func (q *FindQuery) Sort(fields ...string) *FindQuery {
	q.Query.Sort(fields...)
	return q
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
func (q *FindQuery) Hint(indexKeys ...string) *FindQuery {
	q.Query.Hint(indexKeys...)
	return q
}

// Projection limits the information returned from the query.
// Whitelist fields using:     `bson.M{"showThisField": 1}`
// Blacklist fields using:     `bson.M{"someBigStructToHide": 0}`
func (q *FindQuery) Projection(projectionQuery interface{}) *FindQuery {
	q.projection = projectionQuery
	return q
}

// Skip sets the skip value to bypass the given number of entries
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindQuery) Skip(skip int) *FindQuery {
	s64 := int64(skip)
	q.skip = &s64
	return q
}

// Limit sets the max value of responses to return when executing the query.
// A note that when working with larger datasets, it is much more
// performance to compare using collection.FindByDate
func (q *FindQuery) Limit(limit int) *FindQuery {
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
func (q *FindQuery) Timeout(d time.Duration) *FindQuery {
	q.Query.Timeout(d)
	return q
}
