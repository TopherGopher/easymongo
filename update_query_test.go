package easymongo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdate(t *testing.T) {
	setup(t)
	coll := createBatmanArchive(t)
	t.Cleanup(func() {
		coll.Drop()
	})

	t.Run("Update one", func(t *testing.T) {
		is := assert.New(t)
		err := coll.Update(bson.M{"name": "The Joker"}, bson.M{
			"$set": bson.M{"lastEncounter": time.Now()}}).One()
		is.NoError(err, "Could not update the timestamp for a single document")
	})
	t.Run("Update one by ID", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		t.Skipf("TODO: UpdateOneByID()")
	})
	t.Run("Update many", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		t.Skipf("TODO: Update.Many()")
	})
}
