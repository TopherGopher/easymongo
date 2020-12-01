package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	setup(t)
	dbName := "pokemon"
	collName := ""
	coll := conn.Database(dbName).Collection(collName)
	t.Cleanup(func() {
		coll.Drop()
		teardown(t)
	})

	t.Run("Update one", func(t *testing.T) {
		is := assert.New(t)
		_ = is
	})
	t.Run("Update one by ID", func(t *testing.T) {
		is := assert.New(t)
		_ = is
	})
	t.Run("Update many", func(t *testing.T) {
		is := assert.New(t)
		_ = is
	})
}
