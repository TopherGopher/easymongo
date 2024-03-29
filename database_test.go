package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tophergopher/easymongo"
)

func TestGetDatabase(t *testing.T) {
	setup(t)
	var err error
	coll := createBatmanArchive(t)
	dbName := coll.GetDatabase().Name()
	t.Run("GetDatabase", func(t *testing.T) {
		is := assert.New(t)
		db := conn.Database(dbName)
		is.NotNil(db)
	})
	t.Run("DatabaseNames", func(t *testing.T) {
		is := assert.New(t)
		dbNames := conn.DatabaseNames()
		is.NoError(err, "Couldn't find the database names")
		is.Contains(dbNames, dbName, "There should only be one database present")
	})
	t.Run("ListDatabases", func(t *testing.T) {
		is := assert.New(t)
		dbs := conn.ListDatabases()
		is.NoError(err, "Couldn't list the databases")
		found := false
		for _, db := range dbs {
			if db.Name() == dbName {
				found = true
			}
		}
		is.True(found, "Could not find the expected database among the database list")
	})
	t.Run("db.CollectionNames", func(t *testing.T) {
		is := assert.New(t)
		collNames := conn.Database(dbName).CollectionNames()
		is.Len(collNames, 1, "There should only be one collection present")
	})
	t.Run("db.ListCollections", func(t *testing.T) {
		is := assert.New(t)
		colls, err := conn.Database(dbName).ListCollections()
		is.NoError(err)
		is.Len(colls, 1, "There should only be one collection present")
		if len(colls) != 1 || err != nil {
			t.FailNow()
		}
		is.NotNil(colls[0])
		if colls[0] == nil {
			t.FailNow()
		}
		// Now just try a standard Count() operation on that collection object
		is.Equal("enemies", colls[0].Name(), "The collection name on the test appears to be unset. This appears to be an improperly initialized Collection object.")
	})
	t.Run("db.ListCollections", func(t *testing.T) {
		is := assert.New(t)
		colls, err := conn.Database(dbName).ListCollections()
		is.NoError(err, "Couldn't list the collections")
		is.Len(colls, 1, "There should only be one collection present")
	})
	t.Run("easymongo.GetDatabase global", func(t *testing.T) {
		is := assert.New(t)
		colls, err := easymongo.GetDatabase(dbName).ListCollections()
		is.NoError(err, "Couldn't list the collections")
		is.Len(colls, 1, "There should only be one collection present")
	})
	t.Run("db.Drop()", func(t *testing.T) {
		is := assert.New(t)
		dbNamesBefore := conn.DatabaseNames()
		db := conn.Database(dbName)
		err := db.Drop()
		is.NoError(err, "Could not drop the database")
		dbNamesAfter := conn.DatabaseNames()
		dbLengthDiff := len(dbNamesBefore) - len(dbNamesAfter)
		is.Equal(1, dbLengthDiff, "A database should have been dropped")
	})
}
