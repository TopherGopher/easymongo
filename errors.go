package easymongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoErr represents any kind of mongo error, regardless of which
// component of mongo-go-driver threw the error.
type MongoErr struct {
	Message string `json:"message"`
	err     error
}

// NewMongoErr casts any error to a MongoErr. This is intended
// to help extrapolate all mongo errors into a single class.
// Should be used as:
// if v := errors.Is(err, mongo.Error) {
//		merr := NewMongoErr(err)
// }
// TODO: Should we be using this?
func NewMongoErr(err error) (me MongoErr) {
	return MongoErr{
		Message: err.Error(),
		err:     err,
	}
}
func (me MongoErr) Error() string {
	return me.Message
}
func (me MongoErr) Unwrap() error {
	return me.err
}

var (
	// ErrNotImplemented is raised when a function is not yet supported/complete
	// This is mostly used to help track development progress
	ErrNotImplemented = NewMongoErr(errors.New("this feature has not yet been implemented"))
	// ErrTimeoutOccurred denotes a query exceeded the max call time it was allowed
	ErrTimeoutOccurred = NewMongoErr(errors.New("timeout during database transaction"))
	// ErrPointerRequired denotes that the provided result object is being passed by value rather than reference
	ErrPointerRequired = NewMongoErr(errors.New("a pointer is required in order to unpack the resultant value from a query"))
	// ErrNoDocuments denotes no documents were found
	ErrNoDocuments = NewMongoErr(mongo.ErrNoDocuments)
	// ErrWrongType indicates the specified distinct operation did not work. Check the field type that you are attempting to use distinct on.
	ErrWrongType = NewMongoErr(errors.New("the type specified could not be decoded into"))
)
