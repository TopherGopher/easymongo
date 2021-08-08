package easymongo

import (
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find allows a user to execute a standard find() query.
// findOne(), find() and findAnd*() are executed when a user calls:
//      q.One(), q.Many(). q.FindAnd().Replace(), q.FindAnd().Update(), q.FindAnd().Delete()
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
// Sort/Limit/Skip are presently ignored.
func (q *FindQuery) Distinct(fieldName string) (interfaceSlice []interface{}, err error) {
	// opts := q.findDistinctOptions()
	// mongoColl := q.collection.mongoColl
	// ctx, cancelFunc := q.getContext()
	// defer cancelFunc()
	// interfaceSlice, err = mongoColl.Distinct(ctx, fieldName, q.filter, opts)
	// err = q.collection.handleErr(err)
	pipeline := []bson.M{
		0: {
			"$match": q.filter,
		},
	}

	if q.sortFields != nil {
		pipeline = append(pipeline, bson.M{
			"$sort": *q.sortFields,
		})
	}

	if q.skip != nil && *q.skip > 0 {
		pipeline = append(pipeline, bson.M{
			"$skip": *q.skip,
		})
	}
	if q.limit != nil && *q.limit >= 0 {
		pipeline = append(pipeline, bson.M{
			"$limit": *q.limit,
		})
	}
	pipeline = append(pipeline, []bson.M{
		{
			"$group": bson.M{
				"_id": 0,
				"distinctValues": bson.M{
					"$addToSet": "$" + fieldName,
				},
			},
		},
		// {
		// 	"$project": bson.M{
		// 		"_id": "$distinctValues",
		// 		// 		"arrayAsObj": bson.M{
		// 		// 			"$arrayToObject": "$distinctValues", // Making array an object
		// 		// 		},
		// 	},
		// },
		// {
		// 	"$replaceRoot": bson.M{
		// 		"newRoot": "$distinctValues",
		// 	},
		// },
	}...)

	// { $replaceRoot: { newRoot: { $ifNull: [ "$name", { _id: "$_id", missingName: true} ] } } }
	// distinctDocuments := []distinctDocument{}
	d := distinctDocument{}
	err = q.collection.Aggregate(pipeline).One(&d)
	// if err != nil {
	// 	return interfaceSlice, err
	// } else if len(distinctDocuments) == 1 && len(distinctDocuments[0].distinctValues) > 0 {
	// 	return distinctDocuments[0].distinctValues, nil
	// }
	return d.DistinctValues, err
}

type distinctDocument struct {
	DistinctValues []interface{} `bson:"distinctValues"`
}

// DistinctInts returns a list of distinct integers for a given field.
func (q *FindQuery) DistinctInts(fieldName string) (intSlice []int, err error) {
	iSlice, err := q.Distinct(fieldName)
	if err != nil {
		return intSlice, err
	}
	intSlice = make([]int, len(iSlice))
	for i, iFace := range iSlice {
		switch val := iFace.(type) {
		case int:
			intSlice[i] = val
		case int32:
			intSlice[i] = int(val)
		case int64:
			intSlice[i] = int(val)
		default:
			return intSlice, fmt.Errorf("the field '%s' had values that could not be coerced to ints - raw value type: %T example: %v", fieldName, val, val)
		}
	}
	// TODO: For whatever reason, sort doesn't seem to work in the pipeline aggregation when dealing with ints
	sort.Ints(intSlice)
	return intSlice, nil
}

// DistinctFloat64s returns a list of distinct float64s for a given field.
func (q *FindQuery) DistinctFloat64s(fieldName string) (floatSlice []float64, err error) {
	iSlice, err := q.Distinct(fieldName)
	if err != nil {
		return floatSlice, err
	}
	floatSlice = make([]float64, len(iSlice))
	for i, iFace := range iSlice {
		switch val := iFace.(type) {
		case float32:
			floatSlice[i] = float64(val)
		case float64:
			floatSlice[i] = val
		default:
			return floatSlice, fmt.Errorf("the field '%s' had values that could not be coerced to float64 - raw value type: %T example: %v", fieldName, val, val)
		}
	}
	// TODO: For whatever reason, sort doesn't seem to work in the pipeline aggregation when dealing with numbers
	// Need to be conditional about this sort
	sort.Float64s(floatSlice)
	return floatSlice, nil
}

// DistinctStrings returns a distinct list of strings using the provided query/field name.
// Sorting/Limiting/Skipping are supported, with the caveat that only the Skip/Sort/Limit options are supported
//     coll.Find().Sort(fieldName).Limit(2).Skip(1).DistinctStrings(fieldName)
func (q *FindQuery) DistinctStrings(fieldName string) (stringSlice []string, err error) {
	// iSlice, err := q.Distinct(fieldName)
	// if err != nil {
	// 	return stringSlice, err
	// }
	opts := q.findDistinctOptions()
	mongoColl := q.collection.mongoColl
	ctx, cancelFunc := q.getContext()
	defer cancelFunc()
	iSlice, err := mongoColl.Distinct(ctx, fieldName, q.filter, opts)
	if err = q.collection.handleErr(err); err != nil {
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
		for _, sortFieldE := range *q.sortFields {
			reverseSort := false
			if sortFieldE.Key == fieldName {
				sort.Strings(stringSlice)
				if !reverseSort {
					// If the sort wasn't specified in the provided fieldName, use the sort value from the bson.E value
					reverseSort = sortFieldE.Value.(int) == -1
				}

				if reverseSort {
					sort.Sort(sort.Reverse(sort.StringSlice(stringSlice)))
				}
				break
			}
		}
	}

	if q.skip != nil && *q.skip > 0 {
		if len(stringSlice) > int(*q.skip) {
			// The slice is longer than the skip
			// Walk the slice to the specified skip
			stringSlice = stringSlice[int(*q.skip):]
		} else {
			// The slice is greater or equal to the skip - return nothing
			stringSlice = []string{}
		}
	}

	if q.limit != nil && *q.limit >= 0 && len(stringSlice) > int(*q.limit) {
		// Trim the slice using limit
		stringSlice = stringSlice[0:int(*q.limit)]
	}
	return stringSlice, err
}
