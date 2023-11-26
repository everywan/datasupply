package local

import (
	"context"
	"reflect"
	"testing"
)

func TestForward(t *testing.T) {
	key := "payload"
	value := []byte{1, 2}
	result, err := Forward(context.Background(), []string{key}, value)
	if err != nil {
		t.Fatal(err)
	}
	for _key, _value := range result {
		if _key != key {
			t.Fatal("forward output key not match input")
		}
		if !reflect.DeepEqual(_value, value) {
			t.Fatal("forward output value not match input")
		}
	}
}
