package easymongo

import (
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find allows a user to execute a standard find() query.
// findOne(), find() and findAnd*() is run when a user calls:
//      q.One(), q.Many(). q.FindAnd().Replace(), q.FindAnd().Update() and q.FindAnd().Delete()
// TODO: Consider using bsoncore.Doc rather than interface?
func (c *Collection) Find(filter interface{}) (q *FindQuery) {
	q = &FindQuery{
		Query: c.query(filter),
	}
	return q
}

// FindQuery is a helper for finding and counting documents
type FindQuery struct {
	*Query
	skip                *int64
	limit               *int64
	allowDiskUse        *bool
	projection          interface{}
	allowPartialResults *bool
	batchSize           *int32
	maxTime             *time.Duration
	// cursorType          *CursorType
	// max             interface{}
	// maxAwaitTime    *time.Duration
	// min             interface{}
	noCursorTimeout *bool
	oplogReplay     *bool
	returnKey       *bool
	showRecordID    *bool
	snapshot        *bool
	sort            interface{}
}

// TODO: FindQuery.CursorType helpers
// func (q *FindQuery) CursorType(x *CursorType) *FindQuery {
// 	return q
// }

// func (q *FindQuery) TODO: Max(x interface{}) *FindQuery {
// 	return q
// }
// func (q *FindQuery) TODO: MaxAwaitTime(x *time.Duration) *FindQuery {
// 	return q
// }
// func (q *FindQuery) TODO: Min(x interface{}) *FindQuery {
// 	return q
// }
// func (q *FindQuery) TODO: NoCursorTimeout(x *bool) *FindQuery {
// 	return q
// }
// func (q *FindQuery) TODO: OplogReplay(x *bool) *FindQuery {
// 	return q
// }

// func (q *FindQuery) TODO: ReturnKey(x *bool) *FindQuery {
// 	return q
// }
// func (q *FindQuery) TODO: ShowRecordID(x *bool) *FindQuery {
// 	return q
// }

// func (q *FindQuery) TODO: Snapshot(x *bool) *FindQuery {
// 	return q
// }

// All executes the specified query using find() and unmarshals
// the result into the provided interface. Ensure interface{} is either
// a slice or a pointer to a slice.
func (q *FindQuery) All(results interface{}) error {
	// TODO: Check kind to make sure results is a slice or map
	opts := q.findOptions()
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	cursor, err := q.collection.mongoColl.Find(ctx, q.filter, opts)
	if err != nil {
		return err
	}

	// TODO: Inject ErrNotFound if option specified
	err = cursor.All(ctx, results)
	err = q.collection.handleErr(err)
	return err
}

// findOneOptions generates the native mongo driver FindOneOptions from the FindQuery
func (q *FindQuery) findOneOptions() *options.FindOneOptions {
	o := &options.FindOneOptions{
		AllowPartialResults: q.allowPartialResults,
		BatchSize:           q.batchSize,
		Collation:           q.collation,
		Comment:             q.comment,
		MaxTime:             q.timeout,
		// Projection:          q.projection,
		Skip: q.skip,
	}
	if q.hintIndices != nil {
		o.Hint = *q.hintIndices
	}
	if q.sortFields != nil {
		o.Sort = *q.sortFields
	}
	return o
}

// One consumes the specified query and marshals the result
// into the provided interface.
func (q *FindQuery) One(result interface{}) (err error) {
	if !interfaceIsUnpackable(result) {
		return ErrPointerRequired
	}
	opts := q.findOneOptions()
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	err = q.collection.mongoColl.FindOne(ctx, q.filter, opts).Decode(result)
	err = q.collection.handleErr(err)
	if err != nil {
		return err
	}

	// Depending on ErrNotFound setting, consider unsetting ErrNotFound to make consistent experience
	// FindOne.Decode() is the only mongo-go-driver function that returns ErrNotFound when
	// no result was found.
	return err
}

// findOptions generates the native mongo driver FindOptions from the FindQuery
func (q *FindQuery) findOptions() *options.FindOptions {
	o := &options.FindOptions{
		AllowDiskUse:        q.allowDiskUse,
		AllowPartialResults: q.allowPartialResults,
		BatchSize:           q.batchSize,
		Collation:           q.collation,
		Comment:             q.comment,
		Limit:               q.limit,
		MaxTime:             q.timeout,
		Projection:          q.projection,
		Skip:                q.skip,
		// CursorType:          x,
		// Max:                 x,
		// MaxAwaitTime:        x,
		// Min:                 x,
		// NoCursorTimeout:     x,
		// OplogReplay:         x,
		// ReturnKey:           x,
		// ShowRecordID: x,
		// Snapshot:     x,
	}
	if q.hintIndices != nil {
		o.Hint = *q.hintIndices
	}
	if q.sortFields != nil {
		o.Sort = *q.sortFields
	}
	return o
}

// Cursor results the mongo.Cursor. This is useful when working with large numbers of results.
// Alternatively, consider calling collection.Find().One() or collection.Find().All().
func (q *FindQuery) Cursor() (*mongo.Cursor, error) {
	// TODO: Check kind ot make sure it's a slice or map
	opts := q.findOptions()

	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	return q.collection.mongoColl.Find(ctx, q.filter, opts)
}

// countOptions generates the native mongo driver CountOptions from the FindQuery
func (q *FindQuery) countOptions() *options.CountOptions {
	o := &options.CountOptions{
		Limit:     q.limit,
		Skip:      q.skip,
		MaxTime:   q.timeout,
		Collation: q.collation,
	}
	if q.hintIndices != nil {
		o.Hint = *q.hintIndices
	}
	return o
}

// Count counts the number of documents using the specified query
func (q *FindQuery) Count() (int, error) {
	opts := q.countOptions()
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	count, err := mongoColl.CountDocuments(ctx, q.filter, opts)
	err = q.collection.handleErr(err)
	return int(count), err
}

func (q *FindQuery) findDistinctOptions() *options.DistinctOptions {
	return &options.DistinctOptions{
		Collation: q.Query.collation,
		MaxTime:   q.Query.timeout,
	}
}

// Distinct returns an array of the distinct elements in the provided fieldName.
// A note that interfaceSlice does not contain the full document but rather just the
// value from the provided field.
func (q *FindQuery) Distinct(fieldName string) (interfaceSlice []interface{}, err error) {
	opts := q.findDistinctOptions()
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	interfaceSlice, err = mongoColl.Distinct(ctx, fieldName, q.filter, opts)
	err = q.collection.handleErr(err)
	return interfaceSlice, err
}

// DistinctStrings returns a distinct list of strings using the provided query/field name.
func (q *FindQuery) DistinctStrings(fieldName string) (stringSlice []string, err error) {
	iSlice, err := q.Distinct(fieldName)
	if err != nil {
		return stringSlice, err
	}
	stringSlice = make([]string, len(iSlice))
	for i, iFace := range iSlice {
		if val, ok := iFace.(string); ok {
			stringSlice[i] = val
		} else {
			return stringSlice, ErrWrongType
		}
	}
	if q.sortFields != nil {
		m := q.sortFields.Map()
		if _, found := m[fieldName]; found {
			// A sort was specified on this field
			sort.Strings(stringSlice)
		}
		if _, found := m["-"+fieldName]; found {
			// A reverse sort was specified on this field
			sort.Sort(sort.Reverse(sort.StringSlice(stringSlice)))
		}
	}
	return stringSlice, err
}
