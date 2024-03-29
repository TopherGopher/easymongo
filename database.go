package easymongo

import (
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
	return GetCurrentConnection().Database(dbName)
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

// CollectionNames returns the names of the collections as strings.
// If no collections could be found, then an empty list is returned.
func (db *Database) CollectionNames() []string {
	ctx, cancelFunc := db.connection.operationCtx()
	defer cancelFunc()
	opts := options.ListCollections().SetNameOnly(true)
	collectionNames, err := db.mongoDB.ListCollectionNames(ctx, bson.M{}, opts)
	if err != nil {
		return []string{}
	}
	return collectionNames
}

// ListCollections returns a list of Collection objects that can be queried against.
// If you just need the collection names as strings, use db.CollectionNames() instead
func (db *Database) ListCollections() ([]*Collection, error) {
	collectionNames := db.CollectionNames()
	if len(collectionNames) == 0 {
		return []*Collection{}, mongo.ErrNoDocuments
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
func (db *Database) Run(cmd interface{}, result interface{}) error {
	ctx, cancelFunc := db.connection.defaultQueryCtx()
	defer cancelFunc()
	return db.mongoDB.RunCommand(ctx, cmd).Decode(result)
}

// TODO: DB.Login
// func (db *Database) Login(user, pass string) error {
// 	return
// }
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
	ctx, cancel := db.connection.operationCtx()
	defer cancel()
	return db.mongoDB.Drop(ctx)
}

// Name returns the name of the database
func (db *Database) Name() string {
	return db.dbName
}

// TODO: DB.FindRef
// func (db *Database) FindRef(ref *DBRef) *Query {return }
