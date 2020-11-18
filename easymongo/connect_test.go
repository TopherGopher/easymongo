package easymongo_test

import (
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	type args struct {
		mongoURI string
	}
	tests := []struct {
		name string
		args args
		want *Connection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Connect(tt.args.mongoURI); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}
