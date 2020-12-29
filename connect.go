package easymongo

import (
	"context"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConnectOptions holds helpers for configuring a new mongo connection.
type MongoConnectOptions struct {
	mongoURI                string
	connectTimeout          *time.Duration
	defaultQueryTimeout     *time.Duration
	defaultOperationTimeout *time.Duration
	// if a nil slice should encode as null instead of an empty array type, this should be true
	nilSlicesAreNull *bool
	// This is used as the writeconcern w value which requests acknowledgement that write operations propagate to the specified number of mongod instances
	numWritesForConsensus      *int
	runHealthCheckOnConnection bool
	connectionFlag             *ConnectionFlag
}

// clientOptions returns the standard options.ClientOptions that mongo driver is looking for
func (mopts MongoConnectOptions) clientOptions() *options.ClientOptions {
	var opts *options.ClientOptions
	if mopts.connectionFlag != nil {
		opts = mopts.connectionFlag.mongoDriverClientOptions().ApplyURI(mopts.mongoURI)
	} else {
		opts = DefaultAnywhere.mongoDriverClientOptions().ApplyURI(mopts.mongoURI)
	}

	registry := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(registry)
	if mopts.nilSlicesAreNull != nil && *mopts.nilSlicesAreNull {
		// The mongo-driver will set unintialized slices to a null type rather than array type by default.
		// If a user specifies that they desire this behavior, this is a no-op.
	} else {
		// Typical use-case for easymongo - nil slices are saved as array types in mongo to make queries
		// involving slice mutation less prone to error
		nilSliceCodec := bsoncodec.NewSliceCodec(bsonoptions.SliceCodec().SetEncodeNilAsEmpty(true))
		registry.RegisterDefaultEncoder(reflect.Slice, nilSliceCodec)
	}

	// opts.SetRegistry(registry.Build())
	if mopts.connectTimeout != nil {
		// Limit how long to wait to find an available server before erroring (default 30 seconds)
		opts.SetServerSelectionTimeout(*mopts.connectTimeout)
		// Limit how long to wait before a connection is established (default 30 seconds)
		opts.SetConnectTimeout(*mopts.connectTimeout)
	}
	if mopts.numWritesForConsensus != nil {
		opts.SetWriteConcern(writeconcern.New(writeconcern.W(*mopts.numWritesForConsensus)))
	}

	// TODO: TLS config
	// tlsConfig := &tls.Config{}
	// if opts.SSLAllowInvalidCert || opts.SSLAllowInvalidHost {
	// 	tlsConfig.InsecureSkipVerify = true
	// }
	// if opts.SSLPEMKeyFile != "" {
	// 	_, err := addClientCertFromFile(tlsConfig, opts.SSLPEMKeyFile, opts.SSLPEMKeyPassword)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error configuring client, can't load client certificate: %v", err)
	// 	}
	// }
	// if opts.SSLCAFile != "" {
	// 	if err := addCACertFromFile(tlsConfig, opts.SSLCAFile); err != nil {
	// 		return nil, fmt.Errorf("error configuring client, can't load CA file: %v", err)
	// 	}
	// }
	// mopts.SetTLSConfig(tlsConfig)

	return opts
}

// // addCACertFromFile adds a root CA certificate to the configuration given a path
// // to the containing file.
// func addCACertFromFile(cfg *tls.Config, file string) error {
// 	data, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return err
// 	}

// 	certBytes, err := loadCert(data)
// 	if err != nil {
// 		return err
// 	}

// 	cert, err := x509.ParseCertificate(certBytes)
// 	if err != nil {
// 		return err
// 	}

// 	if cfg.RootCAs == nil {
// 		cfg.RootCAs = x509.NewCertPool()
// 	}

// 	cfg.RootCAs.AddCert(cert)

// 	return nil
// }

// Connect connects to the given mongo URI.
// Connect wraps ConnectWithOptions. If you are using just a mongoUri for connection, this should be all you
// need. However, if you need to configure additional options, it is recommened to instead use ConnectWithOptions.
// If a connection does not suceed when using Connect, then a panic occurs.
func Connect(mongoURI string) *Connection {
	connection, err := ConnectWith(mongoURI).Connect()
	if err != nil {
		panic(err)
	}
	return connection
}

// ConnectionBuilder is for specifying options when connecting to a DB
type ConnectionBuilder struct {
	connection Connection
}

// ConnectWith allows one to start building a Connection with Options. Call .Connect() at the end to establish the connection.
// e.g. err = easymongo.ConnectWith(mongoURI).DefaultTimeout().Connect()
func ConnectWith(mongoURI string) *ConnectionBuilder {
	return &ConnectionBuilder{
		Connection{
			mongoOptions: MongoConnectOptions{
				mongoURI: mongoURI,
			},
		},
	}
}

// Flags can be used to set one or more connection flags. Consider using Default* options, or use bitwise '|' to specify multiple options
// e.g.: ConnectWith(mongoURI).Flags(ReadConcernMajority | ReadPreferenceNearest | WriteConcernMajority)
func (cb *ConnectionBuilder) Flags(flags ConnectionFlag) *ConnectionBuilder {
	cb.connection.mongoOptions.connectionFlag = &flags
	return cb
}

// DefaultQueryTimeout allows you to specify a timeout used for query operations.
func (cb *ConnectionBuilder) DefaultQueryTimeout(timeout time.Duration) *ConnectionBuilder {
	cb.connection.mongoOptions.defaultQueryTimeout = &timeout
	return cb
}

// DefaultOperationTimeout allows you to specify a timeout used both for the top-level connection and
// for any subsequent queries/operations to the database (unless overridden).
func (cb *ConnectionBuilder) DefaultOperationTimeout(timeout time.Duration) *ConnectionBuilder {
	cb.connection.mongoOptions.defaultOperationTimeout = &timeout
	return cb
}

// ConnectTimeout allows one to specify the initial timeout when connecting to a database
func (cb *ConnectionBuilder) ConnectTimeout(timeout time.Duration) *ConnectionBuilder {
	cb.connection.mongoOptions.connectTimeout = &timeout
	return cb
}

// FromMongoDriverClient accepts an initialized mongo.Client.
// This is useful if you want the power of standing up your own mongo.Client connection externally.
func (cb *ConnectionBuilder) FromMongoDriverClient(client *mongo.Client) *Connection {
	cb.connection.client = client
	conn := &cb.connection
	// conn := &Connection{
	// 	client: client,
	// }
	setGlobalConnection(conn)
	return conn
}

// Connect performs the actual connection to the DB. A note that calling this function has the
// side-effect of setting the global cached connection to this value. If you are not using the global
// connection value and instead using the value explicitly returned from this function, then disregard this side-effect.
// If a connection does not succeed, then an error is returned.
func (cb *ConnectionBuilder) Connect() (*Connection, error) {
	opts := cb.connection.mongoOptions.clientOptions()
	if cb.connection.client == nil {
		client, err := mongo.NewClient(opts)
		if err != nil {
			return nil, err
		}
		cb.connection.client = client
	}

	ctx, cancel := cb.connection.operationCtx()
	defer cancel()
	if err := cb.connection.client.Connect(ctx); err != nil {
		return nil, err
	}
	setGlobalConnection(&cb.connection)
	return &cb.connection, nil
}

// Connection represents a connection to the mongo cluster/instance.
type Connection struct {
	mongoOptions MongoConnectOptions
	client       *mongo.Client
}

// MongoURI returns the URI that the mongo instance is connected to
func (conn *Connection) MongoURI() string {
	return conn.mongoOptions.mongoURI
}

// Ping attempts to ping the mongo instance
func (conn *Connection) Ping() (err error) {
	ctx, cancel := conn.operationCtx()
	defer cancel()

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

type ConnectionFlag uint8

const (
	ReadConcernAvailable ConnectionFlag = iota + 1
	ReadConcernLinearizable
	ReadConcernLocal
	ReadConcernMajority
	ReadConcernSnapshot

	ReadPreferenceNearest
	// ReadPreferencePrimary limits the read to the primary node. If the primary isn't available, the query will error.
	ReadPreferencePrimary
	// ReadPreferencePrimaryPreferred prefers reads from primary, but will fall back
	ReadPreferencePrimaryPreferred
	ReadPreferenceWriteConcern
	ReadPreferenceSecondary
	ReadPreferenceSecondaryPreferred

	WriteConcernJournal
	WriteConcernW1
	WriteConcernW2
	WriteConcernW3
	// WriteConcernMajority waits for a majority of
	WriteConcernMajority
)

const (
	// DefaultPrimary attempts to read and write to the primary node. Secondary will be read from if primary isn't available.
	DefaultPrimary = ReadConcernLocal | ReadPreferencePrimaryPreferred | WriteConcernW1
	// DefaultSecondary reads by default from a secondary node (if available). It will use majority consensus when reading to determine what data to return.
	DefaultSecondary = ReadConcernMajority | ReadPreferenceSecondaryPreferred | WriteConcernW1
	// DefaultAnywhere connects to the first available node (primary or secondary) for reading. It uses majority both for write confirmations and while waiting for reads.
	DefaultAnywhere = ReadConcernMajority | ReadPreferenceNearest | WriteConcernMajority
)

// mongoDriverDatabaseOptions returns the associated options for the provided connectFlag(s)
func (connectFlag ConnectionFlag) mongoDriverDatabaseOptions() *options.DatabaseOptions {
	opts := options.Database()
	wOpts := make([]writeconcern.Option, 0)

	switch {
	case connectFlag&ReadConcernAvailable == 1:
		opts.SetReadConcern(readconcern.Available())
	case connectFlag&ReadConcernLinearizable == 1:
		opts.SetReadConcern(readconcern.Linearizable())
	case connectFlag&ReadConcernLocal == 1:
		opts.SetReadConcern(readconcern.Local())
	case connectFlag&ReadConcernMajority == 1:
		opts.SetReadConcern(readconcern.Majority())
	case connectFlag&ReadConcernSnapshot == 1:
		opts.SetReadConcern(readconcern.Snapshot())
	}
	switch {
	case connectFlag&ReadPreferenceNearest == 1:
		opts.SetReadPreference(readpref.Nearest())
	case connectFlag&ReadPreferencePrimary == 1:
		opts.SetReadPreference(readpref.Primary())
	case connectFlag&ReadPreferencePrimaryPreferred == 1:
		opts.SetReadPreference(readpref.PrimaryPreferred())
	case connectFlag&ReadPreferenceSecondary == 1:
		opts.SetReadPreference(readpref.Secondary())
	case connectFlag&ReadPreferenceSecondaryPreferred == 1:
		opts.SetReadPreference(readpref.SecondaryPreferred())
	}

	if connectFlag&WriteConcernJournal == 1 {
		wOpts = append(wOpts, writeconcern.J(true))
	}
	if connectFlag&WriteConcernW1 == 1 {
		wOpts = append(wOpts, writeconcern.W(1))
	}
	if connectFlag&WriteConcernW2 == 1 {
		wOpts = append(wOpts, writeconcern.W(2))
	}
	if connectFlag&WriteConcernW3 == 1 {
		wOpts = append(wOpts, writeconcern.W(3))
	}
	if connectFlag&WriteConcernMajority == 1 {
		wOpts = append(wOpts, writeconcern.WMajority())
	}
	if len(wOpts) != 0 {
		opts.SetWriteConcern(writeconcern.New(wOpts...))
	}
	return opts
}

// mongoDriverClientOptions returns the associated options for the provided connectFlag(s)
func (connectFlag ConnectionFlag) mongoDriverClientOptions() *options.ClientOptions {
	opts := options.Client()
	wOpts := make([]writeconcern.Option, 0)

	switch {
	case connectFlag&ReadConcernAvailable == 1:
		opts.SetReadConcern(readconcern.Available())
	case connectFlag&ReadConcernLinearizable == 1:
		opts.SetReadConcern(readconcern.Linearizable())
	case connectFlag&ReadConcernLocal == 1:
		opts.SetReadConcern(readconcern.Local())
	case connectFlag&ReadConcernMajority == 1:
		opts.SetReadConcern(readconcern.Majority())
	case connectFlag&ReadConcernSnapshot == 1:
		opts.SetReadConcern(readconcern.Snapshot())
	}
	switch {
	case connectFlag&ReadPreferenceNearest == 1:
		opts.SetReadPreference(readpref.Nearest())
	case connectFlag&ReadPreferencePrimary == 1:
		opts.SetReadPreference(readpref.Primary())
	case connectFlag&ReadPreferencePrimaryPreferred == 1:
		opts.SetReadPreference(readpref.PrimaryPreferred())
	case connectFlag&ReadPreferenceSecondary == 1:
		opts.SetReadPreference(readpref.Secondary())
	case connectFlag&ReadPreferenceSecondaryPreferred == 1:
		opts.SetReadPreference(readpref.SecondaryPreferred())
	}

	if connectFlag&WriteConcernJournal == 1 {
		wOpts = append(wOpts, writeconcern.J(true))
	}
	if connectFlag&WriteConcernW1 == 1 {
		wOpts = append(wOpts, writeconcern.W(1))
	}
	if connectFlag&WriteConcernW2 == 1 {
		wOpts = append(wOpts, writeconcern.W(2))
	}
	if connectFlag&WriteConcernW3 == 1 {
		wOpts = append(wOpts, writeconcern.W(3))
	}
	if connectFlag&WriteConcernMajority == 1 {
		wOpts = append(wOpts, writeconcern.WMajority())
	}
	if len(wOpts) != 0 {
		opts.SetWriteConcern(writeconcern.New(wOpts...))
	}
	return opts
}

// DatabaseByConnectionType returns the database object associated with the provided database name
func (conn *Connection) DatabaseByConnectionType(dbName string, connectFlag ConnectionFlag) *Database {
	return &Database{
		connection: conn,
		dbName:     dbName,
		mongoDB:    conn.client.Database(dbName, connectFlag.mongoDriverDatabaseOptions()),
	}
}

// DatabaseNames returns a list of the databases available in the connected cluster as a list of strings.
// If an error occurrent, an empty list is returned.
func (conn *Connection) DatabaseNames() []string {
	opts := options.ListDatabases()
	ctx, cancel := conn.operationCtx()
	defer cancel()
	list, err := conn.client.ListDatabaseNames(ctx, bson.M{}, opts)
	if err != nil {
		// TODO: Should we return the error instead? If we don't change this, we should change CollectionNames()
		list = []string{}
	}
	return list
}

// getDefaultQueryCtx returns a context based on if a default query timeout has been set.
// context.Background() and an empty inlined function are returned if no timeout has been set.
func (conn *Connection) defaultQueryCtx() (ctx context.Context, cancel context.CancelFunc) {
	return GetTimeoutCtx(conn.mongoOptions.defaultQueryTimeout)
}

// operationCtx returns a context based on if a default operation timeout has been set.
// context.Background() and an empty inlined function are returned if no timeout has been set.
func (conn *Connection) operationCtx() (ctx context.Context, cancel context.CancelFunc) {
	return GetTimeoutCtx(conn.mongoOptions.defaultOperationTimeout)
}

// GetTimeoutCtx returns a context based on if a timeout has been specified. If no timeout
// was specified, then context.Background() is returned.
func GetTimeoutCtx(timeout *time.Duration) (ctx context.Context, cancel context.CancelFunc) {
	ctx = context.Background()
	// Make cancel a no-op function by default to avoid possible nil function calls
	// Empty inlined functions end up no-oped by compiler
	cancel = func() {}
	if timeout != nil {
		ctx, cancel = context.WithTimeout(context.Background(), *timeout)
	}
	return ctx, cancel
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
		panic("Connect() or ConnectWith() must be called prior to GetCurrentConnection()")
	}
	return globalConnection
}
