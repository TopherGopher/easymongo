package easymongo

import "reflect"

// interfaceIsZero returns true when an interface is either 0 or nil
func interfaceIsZero(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// interfaceIsUnpackable is a helper for determining whether we will be able
// to unpack a result into the given interface. Since Slices and Maps are natively
// pointers, it is acceptable for a Kind to be a Slice and still have things work out.
func interfaceIsUnpackable(x interface{}) bool {
	val := reflect.ValueOf(x)
	kind := val.Kind()
	switch kind {
	case reflect.Ptr:
		return true
	case reflect.Slice:
		return true
	case reflect.Map:
		return true
	}
	return false
}
