package easymongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is helper for representing a database in a cluster.
type Database struct {
	connection *Connection
	dbName     string
	mongoDB    *mongo.Database
}

// GetDatabase returns a database object for the named database using the most recently connected to
// mongo instance/cluster.
func GetDatabase(dbName string) *Database {
	return GetCurrentConnection().GetDatabase(dbName)
}

// C returns a Collection object that can be used to run queries on and create/drop indices
// C wraps Collection() for users who want to use short-hand to do queries.
// It is interchangeable with a call to db.Collection().
func (db *Database) C(name string) *Collection {
	return db.Collection(name)
}

// Collection returns a Collection object that can be used to run queries on and create/drop indices
// It is interchangeable with C().
func (db *Database) Collection(name string) *Collection {
	opts := options.Collection()
	return &Collection{
		database:       db,
		collectionName: name,
		mongoColl:      db.mongoDB.Collection(name, opts),
	}
}

// DefaultCtx returns the appropriate context using the default timeout specified at conneciton time.
func (db *Database) DefaultCtx() (context.Context, context.CancelFunc) {
	return db.connection.GetDefaultTimeoutCtx()
}

// CollectionNames returns the names of the collections as strings.
func (db *Database) CollectionNames() (collectionNames []string, err error) {
	ctx, cancelFunc := db.DefaultCtx()
	defer cancelFunc()
	opts := options.ListCollections().SetNameOnly(true)
	collectionNames, err = db.mongoDB.ListCollectionNames(ctx, bson.M{}, opts)
	return collectionNames, err
}

// ListCollections returns a list of Collection objects that can be actioned against.
func (db *Database) ListCollections() ([]*Collection, error) {
	collectionNames, err := db.CollectionNames()
	if collectionNames == nil || err != nil {
		return []*Collection{}, err
	}
	colls := make([]*Collection, len(collectionNames))
	for i, collName := range collectionNames {
		colls[i] = db.Collection(collName)
	}
	return colls, nil
}

// TODO: DB.With
// func (db *Database) With(s *Session) *Database {return }
// TODO: DB.GridFS
// func (db *Database) GridFS(prefix string) *GridFS {return }
// TODO: DB.Run
// func (db *Database) Run(cmd interface{}, result interface{}) error {return }
// TODO: DB.Login
// func (db *Database) Login(user, pass string) error {return }
// TODO: DB.Logout
// func (db *Database) Logout() {return }
// TODO: DB.UpsertUser
// func (db *Database) UpsertUser(user *User) error {return }
// TODO: DB.runUserCmd
// func (db *Database) runUserCmd(cmdName string, user *User) error {return }
// TODO: DB.AddUser
// func (db *Database) AddUser(username, password string, readOnly bool) error {return }
// TODO: DB.RemoveUser
// func (db *Database) RemoveUser(user string) error {return }

// Drop drops a database from a mongo instance. Use with caution.
func (db *Database) Drop() error {
	ctx, cancel := db.DefaultCtx()
	defer cancel()
	return db.mongoDB.Drop(ctx)
}

// TODO: DB.FindRef
// func (db *Database) FindRef(ref *DBRef) *Query {return }
