package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDistinct(t *testing.T) {
	is := assert.New(t)
	setup(t)
	c := createBatmanArchive(t)
	is.NoError(c.Connection().EnableDebug(), "Could not enable connection debug")
	t.Run("Distinct fields", func(t *testing.T) {
		is := assert.New(t)
		nameISlice, err := c.Find(bson.M{}).Distinct("name")
		is.NoError(err, "Finding distinct interface slice failed")
		is.GreaterOrEqual(len(nameISlice), 3, "There should be several items in this list")
		foundName := false
		for i, nameIFace := range nameISlice {
			// Ensure everything can be coerced to the correct type
			name, isString := nameIFace.(string)
			is.True(isString, "Could not coerce returned value ('%v') to string at %d.", nameIFace, i)
			if name == "Edward Nigma" {
				foundName = true
			}
		}
		is.Contains(nameISlice, "Edward Nigma", "Could not find expected value in distinct results")
		is.True(foundName, "Could not find expected value in distinct results")
	})

	t.Run("Distinct strings", func(t *testing.T) {
		is := assert.New(t)
		names, err := c.Find(bson.M{}).DistinctStrings("name")
		is.NoError(err, "Finding distinct strings failed")
		is.GreaterOrEqual(len(names), 3, "There should be several items in this list")
		is.Contains(names, "Edward Nigma", "Could not find expected value in distinct string results")

		names, err = c.Find(bson.M{}).Sort("name").Skip(2).Limit(2).DistinctStrings("name")
		is.NoError(err, "Finding distinct strings failed")
		is.Equal(2, len(names), "There should be several items in this list")
		if len(names) != 2 || err != nil {
			t.FailNow()
		}
		is.Equal("Poison Ivy", names[0], "Something is off in sort/skip/limit for distinct strings")
		is.Equal("Superman", names[1], "Something is off in sort/skip/limit for distinct strings")

		names, err = c.Find(bson.M{}).Sort("-name").Skip(2).Limit(2).DistinctStrings("name")
		is.NoError(err, "Finding distinct strings failed")
		is.Equal(2, len(names), "There should be several items in this list")
		if len(names) != 2 || err != nil {
			t.FailNow()
		}
		is.Equal("Poison Ivy", names[1], "Something is off in sort/skip/limit for distinct strings")
		is.Equal("Superman", names[0], "Something is off in sort/skip/limit for distinct strings")
	})

	t.Run("Distinct ints", func(t *testing.T) {
		is := assert.New(t)
		ints, err := c.Find(bson.M{}).Sort("timesFought").DistinctInts("timesFought")
		is.NoError(err, "Issue finding distinct integers")
		// There should be 4 distinct ints - the numbers 1 through 4
		is.Len(ints, 4, "There should be 4 distinct values")
		if len(ints) != 4 || err != nil {
			t.FailNow()
		}
		for index, val := range ints {
			is.Equal(index+1, val, "There is a mismatch in expected values - perhaps the sort is broken?")
		}
	})

	t.Run("Distinct floats", func(t *testing.T) {
		is := assert.New(t)
		floats, err := c.Find(bson.M{}).Sort("evilness").DistinctFloat64s("evilness")
		is.NoError(err, "Issue finding distinct floats")
		// There should be 4 distinct floats - the numbers 0 through 0.8
		is.Len(floats, 5, "There should be several distinct values")
		if len(floats) != 5 || err != nil {
			t.FailNow()
		}
		for index, val := range floats {
			// Coerce to float32 for testing as the float64s are sometimes a little odd
			// e.g. float64(index=3) * .2 == 0.6000000000000001
			is.Equal(float32(index)*.2, float32(val), "There is a mismatch in expected values - perhaps the sort is broken?")
		}
	})
}
