package easymongo

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConnectOptions holds helpers for configuring a new mongo connection.
type MongoConnectOptions struct {
	connectTimeout *time.Duration
	defaultTimeout *time.Duration
}

// Connect connects to the given mongo URI.
// Connect wraps ConnectWithOptions. If you are using just a mongoUri for connection, this should be all you
// need. However, if you need to configure additional options, it is recommened to instead use ConnectWithOptions.
// If a connection does not suceed when using Connect, then a panic occurs.
func Connect(mongoURI string) *Connection {
	connection, err := ConnectWithOptions(mongoURI, nil)
	if err != nil {
		panic(err)
	}
	return connection
}

// ConnectWithOptions connects to the specified mongo URI. A note that calling this function has the
// side-effect of setting the global cached connection to this value. If you are not using the global
// connection value and instead using the value explicitly returned from this function, then no need to worry about this.
// If a connection does not succeed, then an error is returned.
// TODO: Configure and consume mongoOpts
func ConnectWithOptions(mongoURI string, mongoOpts *MongoConnectOptions) (*Connection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	conn := &Connection{
		client:   client,
		mongoURI: mongoURI,
	}
	if mongoOpts != nil {
		conn.mongoOptions = *mongoOpts
	}
	setGlobalConnection(conn)
	return conn, nil
}

// ConnectUsingMongoClient accepts an initialized mongo.Client and returns an easymongo Connection
// This is useful if you want the power of standing up your own mongo.Client connection.
// mongoURI is used for informational purposes - feel free to ignore it if you don't need it.
func ConnectUsingMongoClient(client *mongo.Client, mongoURI string) *Connection {
	// client.Connect(nil)
	conn := &Connection{
		client:   client,
		mongoURI: mongoURI,
	}
	// TODO: Consider accepting MongoOptions
	// if defaultTimeout != nil {
	// 	conn.mongoOptions.defaultTimeout = defaultTimeout
	// }
	setGlobalConnection(conn)
	return conn
}

// Connection represents a connection to the mongo cluster/instance.
type Connection struct {
	mongoOptions MongoConnectOptions
	client       *mongo.Client
	mongoURI     string
}

// Ping attempts to ping the mongo instance
func (conn *Connection) Ping() (err error) {
	ctx, cancelFunc := conn.GetDefaultTimeoutCtx()
	defer cancelFunc()

	// TODO: Get readpref from singleton
	err = conn.client.Ping(ctx, readpref.PrimaryPreferred())
	if err == nil && ctx != nil && ctx.Err() != nil {
		// A timeout occurrect during the ping
		err = ErrTimeoutOccurred
	}
	return err
}

// MongoDriverClient returns the mongo.Client from mongo-go-driver - to allow
// for direct interaction with the mongo driver for those users searching for
// more fine-grained control.
func (conn *Connection) MongoDriverClient() *mongo.Client {
	return conn.client
}

// Database returns the database object associated with the provided database name
func (conn *Connection) Database(dbName string) *Database {
	opts := options.Database()
	return &Database{
		connection: conn,
		dbName:     dbName,
		mongoDB:    conn.client.Database(dbName, opts),
	}
}

// DatabaseNames returns a list of the databases available in the connected cluster as a list of strings.
// If an error occurrent, an empty list is returned.
func (conn *Connection) DatabaseNames() []string {
	opts := options.ListDatabases()
	ctx, _ := conn.GetDefaultTimeoutCtx()
	list, err := conn.client.ListDatabaseNames(ctx, bson.M{}, opts)
	if err != nil {
		// TODO: Should we return the error instead? If we don't change this, we should change CollectionNames()
		list = []string{}
	}
	return list
}

// GetDefaultTimeoutCtx returns a context based on if a default timeout has been set. If no timeout
// was specified, then context.Background() is returned.
func (conn *Connection) GetDefaultTimeoutCtx() (ctx context.Context, cancelFunc context.CancelFunc) {
	ctx = context.Background()
	// Make cancelFunc a no-op function by default
	cancelFunc = func() {}
	if conn.mongoOptions.defaultTimeout != nil {
		ctx, cancelFunc = context.WithTimeout(
			context.Background(), *conn.mongoOptions.defaultTimeout)
	}
	return ctx, cancelFunc
}

// ListDatabases returns a list of databases available in the connected cluster as objects that can be interacted with.
func (conn *Connection) ListDatabases() (dbList []*Database) {
	dbNames := conn.DatabaseNames()
	dbList = make([]*Database, len(dbNames))
	for i, dbName := range dbNames {
		dbList[i] = conn.Database(dbName)
	}
	return dbList
}

// globalConnection is used to cache the most recent cluster connected to
var globalConnection *Connection

// connectionLock should be used whenever modifications are made to globalConnection
var connectionLock sync.RWMutex

// setGlobalConnection sets the cached global connection to the provided connection value.
func setGlobalConnection(conn *Connection) {
	connectionLock.Lock()
	defer connectionLock.Unlock()
	globalConnection = conn
}

// GetCurrentConnection returns the current connection cached in the global context.
func GetCurrentConnection() *Connection {
	connectionLock.RLock()
	defer connectionLock.RUnlock()
	if globalConnection == nil {
		panic("Connect() or ConnectWithOptions() must be called prior to GetCurrentConnection()")
	}
	return globalConnection
}

// MongoURI returns the URI of the mongo instance the test connection
// is tethered to.
func (conn *Connection) MongoURI() string {
	return conn.mongoURI
}
