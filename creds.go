package easymongo

import "go.mongodb.org/mongo-driver/mongo/options"

// TODO: Something to help with credentials

// WithAuth injects credentials into mongoOptions
func (c *ConnectionBuilder) WithAuth(creds *options.Credential) {
	c.connection.mongoOptions.auth = creds
}
