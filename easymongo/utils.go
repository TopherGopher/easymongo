package easymongo

import (
	"fmt"
	"reflect"
)

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

// interfaceSlice takes any slice and converts it to a slice of interface
// Thanks to https://stackoverflow.com/a/12754757
func interfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() == reflect.Ptr {
		// Dereference the pointer (if necessary)
		s = s.Elem()
	}
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("a non-slice type was provided")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil, nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}
