package easymongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertQuery is a helper for constructing insertion operations
type InsertQuery struct {
	*Query
}

// Insert constructs and returns an InsertQuery object.
// Now run .One() or .Many() using this handle.
func (c *Collection) Insert() *InsertQuery {
	return &InsertQuery{
		Query: c.query(nil),
	}
}

// One is used to insert a single object into a collection
func (iq *InsertQuery) One(objToInsert interface{}) (id *primitive.ObjectID, err error) {
	ctx, cancelFunc := iq.getContext()
	defer cancelFunc()
	// TODO: InsertOne options
	opts := options.InsertOne()
	result, err := iq.collection.mongoColl.InsertOne(ctx, objToInsert, opts)
	if result != nil && !interfaceIsZero(result.InsertedID) {
		if rid, ok := result.InsertedID.(primitive.ObjectID); ok {
			id = &rid
		}
	}
	return id, err
}

// Many inserts a slice into a mongo collection. The ids of the inserted objects are
// returned. If an ID was not available (e.g. it was unset), an empty entry in the slice is returned.
// TODO: What should be the default behavior when an ID isn't returned?
// Note: Many() uses reflect to coerce an interface to an interface slice. This results in
// a minor O(N) performance hit. If inserting large quanities of items and every nanosecond counts,
// cast your slice to a slice interface yourself (which is an O(1) operation), and call ManyFromInterfaceSlice().
func (iq *InsertQuery) Many(objsToInsert interface{}) (ids []*primitive.ObjectID, err error) {
	if interfaceIsZero(objsToInsert) {
		return nil, fmt.Errorf("the value provided to Insert().Many() must be defined")
	}
	iSlice, err := interfaceSlice(objsToInsert)
	if err != nil {
		return nil, err
	}
	return iq.ManyFromInterfaceSlice(iSlice)

}

// ManyFromInterfaceSlice inserts an interface slice into a mongo collection
// If you need to insert large quantities of items and every nanosecond matters,
// then use this function instead of Many.
func (iq *InsertQuery) ManyFromInterfaceSlice(objsToInsert []interface{}) (ids []*primitive.ObjectID, err error) {
	ctx, cancelFunc := iq.getContext()
	defer cancelFunc()
	// TODO: InsertMany options
	opts := options.InsertMany()

	result, err := iq.collection.mongoColl.InsertMany(ctx, objsToInsert, opts)
	if err != nil {
		return ids, err
	}
	ids = make([]*primitive.ObjectID, len(result.InsertedIDs))
	for i, ridIface := range result.InsertedIDs {
		if rid, ok := ridIface.(primitive.ObjectID); ok {
			ids[i] = &rid
		}
	}
	return ids, err
}
