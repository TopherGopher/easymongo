package easymongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Query is used for creating
// To create a query, call a relevant function off of a Collection. The following line
// would return a *Query object:
//   `q := GetDatabase("db_name").C("collection_name").Find(bson.M{"alignment": "Chaotic Neutral"})`
// And that can be consumed as such:
//   `var weirdCharacters []Character`
//   `err = q.Skip(5).Limit(10).Timeout(time.Minute).All(&weirdCharacters)``
// The previous lines would:
// - Access db_name.collection_name
// - Setup the query filter(s) to skip the first 5 entries and limit to 10 entries max - timing out the query after 1 minute.
// - Run the find() query looking for records matching {"alignment": "Chaotic Neutral"}
// - Unpack the results into the weirdCharacters slice
type Query struct {
	// filter holds the query to be executed - typically a bson.M or bson.D value.
	filter      interface{}
	sortFields  bson.D
	hintIndices bson.D
	comment     *string
	collation   *options.Collation
	timeout     *time.Duration
	collection  *Collection
	// many        bool
}

// Query returns an initialized query object from a collection
func (c *Collection) Query(filter interface{}) *Query {
	return &Query{
		filter:     filter,
		collection: c,
		timeout:    c.database.connection.mongoOptions.defaultTimeout,
	}
}

// Sort accepts a list of strings to use as sort fields.
// Prepending a field name with a '-' denotes descending sorting
// e.g. "-name" would sort the "name" field in descending order
func (q *Query) Sort(fields ...string) *Query {
	q.sortFields = make(bson.D, len(fields))
	for i, field := range fields {
		q.sortFields[i] = indexKeyToBsonE(field)
	}
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
func (q *Query) Hint(indexKeys ...string) *Query {
	q.hintIndices = make(bson.D, len(indexKeys))
	for i, indexKey := range indexKeys {
		q.hintIndices[i] = indexKeyToBsonE(indexKey)
	}
	return q
}

// Comment adds a comment to the query - when the query is executed, this
// comment can help with debugging from the logs.
func (q *Query) Comment(comment string) *Query {
	q.comment = &comment
	return q
}

// Timeout uses the provided duration to set a timeout value using
// a context. The timeout clock begins upon query execution (e.g. calling .All()),
// not at time of calling Timeout().
func (q *Query) Timeout(d time.Duration) *Query {
	q.timeout = &d
	return q
}

// Collation allows users to specify language-specific rules for string comparison, such as rules for lettercase and accent marks.
// https://docs.mongodb.com/manual/reference/collation/
// TODO: Create helpers and consts for Collation
func (q *Query) Collation(c *options.Collation) *Query {
	q.collation = c
	return q
}

// SetContext allows one to override the implied context that is created
// at query time and instead will consume this.
// TODO: if we allow this, the context will already be ticking before the
// query is ever executed.
// func (q *Query) SetContext(ctx *context.Context) *Query { return q }

// getContext returns the appropriate context using the Timeout that was specified either by SetTimeout
// at the query level, or by consuming the default top-level timeout (specified at initialization time).
// getContext should be called after the query has been constructed (thus the private specification).
func (q *Query) getContext() (context.Context, context.CancelFunc) {
	if q.timeout != nil {
		return context.WithTimeout(nil, *q.timeout)
	}
	return q.collection.DefaultCtx()
}

//////////////////////////////
// TO BE IMPLEMENTED!!!!
//////////////////////////////

// Explain provides helpful feedback on what the database is planning on executing. This is useful
// to understand why a query is not being performant.
// TODO: Originally this took an interface, but I think this makes more sense to return as a predefined object.
// func (q *Query) Explain(result interface{}) error {
// 	// q.explainResult = result
// 	return ErrNotImplemented
// }

// ExplainResult represents the results from a query documented here - https://docs.mongodb.com/manual/reference/command/explain/#dbcmd.explain
// TODO: Test that this actually unpacks and support the other two types of explain results.
type ExplainResult struct {
	QueryPlanner QueryPlanner `json:"queryPlanner"`
}
type ParsedQuery struct {
}
type InputStage struct {
	Stage       string        `json:"stage"`
	InputStages []interface{} `json:"inputStages"`
}
type Plan struct {
	Stage      string     `json:"stage"`
	InputStage InputStage `json:"inputStage"`
}
type QueryPlanner struct {
	PlannerVersion    string      `json:"plannerVersion"`
	Namespace         string      `json:"namespace"`
	IndexFilterSet    bool        `json:"indexFilterSet"`
	ParsedQuery       ParsedQuery `json:"parsedQuery"`
	QueryHash         string      `json:"queryHash"`
	PlanCacheKey      string      `json:"planCacheKey"`
	OptimizedPipeline bool        `json:"optimizedPipeline"`
	WinningPlan       Plan        `json:"winningPlan"`
	RejectedPlans     []Plan      `json:"rejectedPlans"`
}

// NewQuery is a helper that consumes the global connection to return a query object.
// If you wish to use an explicit Connection object instead, call
// func NewQuery(dbName, collectionName string, query interface{}) {}
// func (q *Query) Batch(n int) *Query {}
// func (q *Query) Prefetch(p float64) *Query {}
// func (q *Query) Select(selector interface{}) *Query {}
// func (q *Query) LogReplay() *Query {}
// func (q *Query) Tail(timeout time.Duration) *Iter {}
// func (q *Query) For(result interface{}, f func() error) error {}
// func (q *Query) Distinct(key string, result interface{}) error {}
// func (q *Query) MapReduce(job *MapReduce, result interface{}) (info *MapReduceInfo, err error) {}
// func (q *Query) Apply(change Change, result interface{}) (info *ChangeInfo, err error) {}
// func (q *Query) SetMaxScan(n int) *Query {}
// func (q *Query) Snapshot() *Query {}
