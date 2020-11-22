package easymongo

import "errors"

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
	ErrNotImplemented = errors.New("this feature has not yet been implemented")
	// ErrTimeoutOccurred denotes a query exceeded the max call time it was allowed
	ErrTimeoutOccurred = NewMongoErr(errors.New("timeout during database transaction"))
	// ErrPointerRequired denotes that the provided result object is being passed by value rather than reference
	ErrPointerRequired = NewMongoErr(errors.New("a pointer is required in order to unpack the resultant value from a query"))
)
