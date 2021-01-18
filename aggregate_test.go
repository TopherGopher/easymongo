package easymongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tophergopher/easymongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAggregate(t *testing.T) {
	setup(t)
	coll := createBatmanArchive(t)

	t.Run("Aggregate().One()", func(t *testing.T) {
		is := assert.New(t)
		// This aggregation query just sums up the number of documents that have notes on them
		// (Which should be 3)
		pipe := []bson.M{
			0: {
				"$match": bson.M{
					"notes": bson.M{"$ne": nil},
				},
			},
			1: {
				"$count": "total",
			},
		}
		var result struct {
			Total int `bson:"total"`
		}
		err := coll.Aggregate(pipe).One(&result)
		is.NoError(err, "could not Aggregate().One()")
		is.Equal(3, result.Total, "the count appears to be incorrect from the aggregation result")
	})
	t.Run("Aggregate().All()", func(t *testing.T) {
		is := assert.New(t)
		pipe := []bson.M{
			0: {
				// Look for all objects
				"$match": bson.M{},
			},
			1: {
				"$group": bson.M{
					"_id": "$deceased",
					"total": bson.M{
						"$sum": 1,
					},
				},
			},
		}
		var result []struct {
			Deceased bool `bson:"_id"`
			Total    int  `bson:"total"`
		}
		err := coll.Aggregate(pipe).All(&result)
		is.NoError(err, "could not Aggregate().All()")
		is.Len(result, 2, "two categories should be returned from the group")
		if len(result) != 2 {
			t.FailNow()
		}
		for _, res := range result {
			if res.Deceased {
				is.Equal(1, res.Total, "there should only be 1 document in this group")
			} else {
				is.GreaterOrEqual(res.Total, 5, "there should be at least 5 documents in this group")
			}
		}
		// output, err := coll.Aggregate(pipe).FetchRawResult()
		// is.NoError(err, "could not fetch the raw output of the query")
		// _ = output
		// is.GreaterOrEqual(len(output), 0, "the raw output appears empty")
		// fmt.Println(output)
	})
	t.Run("Aggregate().One() with options", func(t *testing.T) {
		is := assert.New(t)
		pipe := []bson.M{
			0: {
				"$match": bson.M{
					"notes": bson.M{"$ne": nil},
				},
			},
			1: {
				"$count": "total",
			},
		}
		var result struct {
			Total int `bson:"total"`
		}
		err := coll.Aggregate(pipe).Timeout(time.Minute * 2).One(&result)
		is.NoError(err, "could not Aggregate().One() with a Timeout")
		is.Equal(3, result.Total, "the count appears to be incorrect from the aggregation result")
		ctx := context.Background()
		err = coll.Aggregate(pipe).WithContext(ctx).One(&result)
		is.NoError(err, "could not Aggregate().One() using a Context")
		is.Equal(3, result.Total, "the count appears to be incorrect from the aggregation result")
		// Set a timeout of 0 so we trigger an error
		err = coll.Aggregate(pipe).Timeout(0).One(&result)
		is.Equal(easymongo.ErrTimeoutOccurred, err, "A timeout was expected")
	})
}
