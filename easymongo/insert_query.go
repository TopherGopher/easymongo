package easymongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertQuery struct {
	Query
}

// One is used to insert a single object into a collection
func (iq *InsertQuery) One(objToInsert interface{}) (id *primitive.ObjectID, err error) {
	ctx, cancelFunc := iq.getContext()
	defer cancelFunc()
	// TODO: InsertOne options
	opts := options.InsertOne()
	result, err := iq.collection.mongoColl.InsertOne(ctx, objToInsert, opts)
	if !interfaceIsZero(result.InsertedID) {
		if rid, ok := result.InsertedID.(primitive.ObjectID); ok {
			id = &rid
		}
	}
	return id, err
}

// Many inserts a slice into a mongo collection
func (iq *InsertQuery) Many(objsToInsert interface{}) (ids []*primitive.ObjectID, err error) {
	ctx, cancelFunc := iq.getContext()
	defer cancelFunc()
	// TODO: InsertMany options
	opts := options.InsertMany()

	if interfaceIsZero(objsToInsert) {
		return nil, fmt.Errorf("the value provided to Insert().Many() must be defined")
	}
	iSlice, err := interfaceSlice(objsToInsert)
	if err != nil {
		return nil, err
	}

	result, err := iq.collection.mongoColl.InsertMany(ctx, iSlice, opts)
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
