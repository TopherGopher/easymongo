package easymongo

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

// Iter represents an iterator type. It is used to process the query
// result and unpack it into structs.
type Iter struct {
	cursor   *mongo.Cursor
	query    *FindQuery
	timedOut bool
	isDone   bool
	err      error
	doneLock sync.RWMutex
}

///////////////////////////
// TO BE IMPLEMENTED!!!
/////////////////////////

// Err returns an error should one have occurred while iterating.
// A note that timeout errors are also considered as part of this.
func (iter *Iter) Err() (err error) {
	iter.doneLock.RLock()
	defer iter.doneLock.RUnlock()
	return iter.err
}

// TimedOut returns True if the request timed out.
func (iter *Iter) TimedOut() bool {
	return iter.timedOut
}

// Done is returned when there are no more results to be had from the cursor, whether
// from an error or otherwise
func (iter *Iter) Done() bool {
	iter.doneLock.RLock()
	defer iter.doneLock.RUnlock()
	return iter.isDone
}

// setError is a threadsafe helper for setting an iterator's error.
// The isDone flag is also flipped to True. A note that if setError calls setDone, there's
// potential for a deadlock over the mutex. So...don't do it :-)
func (iter *Iter) setError(err error) {
	iter.doneLock.Lock()
	defer iter.doneLock.Unlock()
	iter.err = err
	iter.isDone = true
}

// setDone flips the isDone flag to true. Subsequent calls to iter.Done() will return True.
func (iter *Iter) setDone() {
	iter.doneLock.Lock()
	defer iter.doneLock.Unlock()
	iter.isDone = true
}

// Next iterates to the next value and unpacks it into the result struct
func (iter *Iter) Next(result interface{}) bool {
	// TODO: Check to make sure this is a pointer
	ctx, cancelFunc := iter.query.getContext()
	defer cancelFunc()
	var err error
	if iter.cursor.Next(ctx) {
		if err = iter.cursor.Decode(result); err != nil {
			iter.setError(err)
			return false
		}
	} else {
		iter.setDone()
	}
	if err = iter.cursor.Err(); err != nil {
		iter.setError(err)
		return false
	}
	if ctx.Err() != nil {
		// TODO: Is it always true that a context error means timeout?
		iter.timedOut = true
		iter.isDone = true
		iter.err = ErrTimeoutOccurred
	}
	// TODO: consume err
	_ = err
	// TODO: Maybe we shouldn't be doing the decoding internally? Try out a few patterns.
	return false
}

// All unpacks all query results into the provided interface.
// The provided interface should be a list or a map.
func (iter *Iter) All(results interface{}) error {
	// TODO: Check to ensure results is a slice or map kind
	ctx, cancelFunc := iter.query.getContext()
	defer cancelFunc()
	err := iter.cursor.All(ctx, results)
	if ctx.Err() != nil {
		iter.timedOut = true
	}
	return err
}

// Close closes the current cursor. This only needs to be called when using `.Next()` as `.All()` automatically
// closes the cursor in mongo-go-driver.
func (iter *Iter) Close() error {
	ctx, cancelFunc := iter.query.getContext()
	defer cancelFunc()
	return iter.cursor.Close(ctx)
}

// TODO: Find out more about For in Iter:
// func (iter *Iter) For(result interface{}, f func() error) (err error) {
// 	return ErrNotImplemented
// }
