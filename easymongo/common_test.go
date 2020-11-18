package easymongo_test

import (
	"github.com/TopherGopher/pkg/mongotest"
)
var conn *mongotest.TestConnection
func setup(t *testing.T) {
	is := assert.New(t)
	var err error
	conn, err = mongotest.NewTestConnection(true)
	is.NoError(err, "Could not stand up test database connection")
}
func teardown() {
	if conn != nil {
		is.NoError(conn.KillMongoContainer(), "Could not tear down test mongo container"
	}
}