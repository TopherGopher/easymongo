package easymongo

import "sync"

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
