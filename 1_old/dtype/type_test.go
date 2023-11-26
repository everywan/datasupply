package dtype

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDType(t *testing.T) {
	dt := String
	dtMarsal, err := dt.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"string"`, string(dtMarsal))

	var dt2 DType
	err = json.Unmarshal(dtMarsal, &dt2)
	assert.NoError(t, err)
	fmt.Println(dt2)
}
