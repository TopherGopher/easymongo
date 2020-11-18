package easymongo_test

import (
	"testing"

	"github.com/TopherGopher/easymongo"
	"github.com/stretchr/testify/assert"
)

func TestGetDatabase(t *testing.T) {
	setup(t)
	defer teardown(t)
	var err error
	dbName, _ := createBatmanArchiveUsingMongoDriver(t)
	t.Run("GetDatabase", func(t *testing.T) {
		is := assert.New(t)
		db := conn.GetDatabase(dbName)
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
		collNames, err := conn.GetDatabase(dbName).CollectionNames()
		is.NoError(err, "Couldn't find the collection names")
		is.Len(collNames, 1, "There should only be one collection present")
	})
	t.Run("db.ListCollections", func(t *testing.T) {
		is := assert.New(t)
		colls, err := conn.GetDatabase(dbName).ListCollections()
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
		db := conn.GetDatabase(dbName)
		err := db.Drop()
		is.NoError(err, "Could not drop the database")
		dbNamesAfter := conn.DatabaseNames()
		dbLengthDiff := len(dbNamesBefore) - len(dbNamesAfter)
		is.Equal(1, dbLengthDiff, "A database should have been dropped")
	})
}
