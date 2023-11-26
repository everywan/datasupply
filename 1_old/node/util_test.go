package node

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamValueCheckNotZero(t *testing.T) {
	t.Parallel()
	nilError := errors.New("is nil")
	zeroError := errors.New("is zero")
	testCases := []struct {
		name        string
		value       interface{}
		expectError error
	}{
		{"int", 100, nil},
		{"string", "fake_value", nil},
		{"array", []int64{1, 2}, nil},
		{"map", map[string]interface{}{"key": 100}, nil},
		{"nil", nil, nilError},
		{"int_0", 0, zeroError},
		{"int64_0", int64(0), zeroError},
		{"string_0", "", zeroError},
		{"array_0", [0]string{}, zeroError},
		{"slice_0", nil, nilError},
	}

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			err := ParamValueCheckNotZero(tcase.value)
			assert.Equal(t, tcase.expectError, err, tcase.name)
		})
	}
}
