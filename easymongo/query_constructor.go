package easymongo

import "go.mongodb.org/mongo-driver/bson"

// UpdateQueryConstructor is a helper for easily crafting update queries.
// TODO: Support multiple calls to each block with unique field names, right now, calling
// .Increment().Increment() would only increment the last field
type UpdateQueryConstructor struct {
	matchQuery interface{}
	setQuery   bson.M
}

// Match sets the Update query filter to the provided interface. If not specified, all
// documents will be updated.
func (qc *UpdateQueryConstructor) Match(query interface{}) *UpdateQueryConstructor {
	qc.matchQuery = query
	return qc
}

// Push pushes the specified object on to an array.
func (qc *UpdateQueryConstructor) Push(arrayFieldName string, objToPush interface{}) *UpdateQueryConstructor {
	qc.setQuery["$push"] = bson.M{arrayFieldName: objToPush}
	return qc
}

// Pull pulls any array values that match the pullCondition query.
func (qc *UpdateQueryConstructor) Pull(arrayFieldName string, pullCondition interface{}) *UpdateQueryConstructor {
	qc.setQuery["$pull"] = bson.M{
		arrayFieldName: pullCondition,
	}
	return qc
}

// PullAll pulls any array values that match the pullCondition query.
// TODO: What is the difference between Pull and PullAll?
func (qc *UpdateQueryConstructor) PullAll(arrayFieldName string, pullCondition interface{}) *UpdateQueryConstructor {
	qc.setQuery["$pullAll"] = bson.M{
		arrayFieldName: pullCondition,
	}
	return qc
}

// AddToSet pushes the provided interface{} object to the specified fieldname
func (qc *UpdateQueryConstructor) AddToSet(arrayFieldName string, objToAdd interface{}) *UpdateQueryConstructor {
	qc.setQuery["$addToSet"] = bson.M{
		arrayFieldName: objToAdd,
	}
	return qc
}

// PopFirst removes the first element from the specified array
func (qc *UpdateQueryConstructor) PopFirst(arrayFieldName string) *UpdateQueryConstructor {
	qc.setQuery["$pop"] = bson.M{
		arrayFieldName: -1,
	}
	return qc
}

// PopLast removes the last element from the specified array
func (qc *UpdateQueryConstructor) PopLast(arrayFieldName string) *UpdateQueryConstructor {
	qc.setQuery["$pop"] = bson.M{
		arrayFieldName: 1,
	}
	return qc
}

// Set sets fieldName to the provided object. If you are looking to replace an entire document, consider
// using collection.Replace() instead.
// e.g. {"$set": {objToSet: objToSet}}
func (qc *UpdateQueryConstructor) Set(fieldName string, objToSet interface{}) *UpdateQueryConstructor {
	qc.setQuery["$set"] = bson.M{
		fieldName: objToSet,
	}
	return qc
}

// Increment increases the specified field value by the provided int.
func (qc *UpdateQueryConstructor) Increment(fieldName string, i int) *UpdateQueryConstructor {
	qc.setQuery["$inc"] = bson.M{
		fieldName: i,
	}
	return qc
}

// Decrement decreases the specified field value by the provided int.
func (qc *UpdateQueryConstructor) Decrement(fieldName string, i int) *UpdateQueryConstructor {
	qc.Increment(fieldName, i*-1)
	return qc
}
