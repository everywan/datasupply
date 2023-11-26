package dtype

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// "github.com/stretchr/testify/assert"

type ConvertTestSuite struct {
	suite.Suite
}

func (suite *ConvertTestSuite) TestToString() {
	testCases := []struct {
		name   string
		value  interface{}
		except string
	}{
		{"int2string", 100, "100"},
		{"bool2string", true, "true"},
		{"bool2string", false, "false"},
	}
	for _, tcase := range testCases {
		suite.Run(tcase.name, func() {
			value, err := Convert(tcase.value, String)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tcase.except, value.(string))
		})
	}
}

func (suite *ConvertTestSuite) TestToInt64() {
	testCases := []struct {
		name   string
		value  interface{}
		except int64
	}{
		{"int2int64", 100, 100},
		{"string2int64", "101", 101},
		{"float2int64", 102.01, 102},
		{"[]byte2int64", []byte("103"), 103},
		{"estring2int64", "1.06059501e+08", 106059501},
	}
	for _, tcase := range testCases {
		suite.Run(tcase.name, func() {
			value, err := Convert(tcase.value, Int64)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tcase.except, value.(int64))
		})
	}
}

func (suite *ConvertTestSuite) TestToFloat64() {
	testCases := []struct {
		name   string
		value  interface{}
		except float64
	}{
		{"int2float64", 100, 100},
		{"string2float64", "101", 101},
		{"float2float64", 102.00, 102},
		{"[]byte2float64", []byte("103"), 103},
		{"estring2float64", "1.06059501e+08", 106059501},
	}
	for _, tcase := range testCases {
		suite.Run(tcase.name, func() {
			value, err := Convert(tcase.value, Float64)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tcase.except, value.(float64))
		})
	}
}

func (suite *ConvertTestSuite) TestToArrayInt64() {
	testCases := []struct {
		name   string
		value  interface{}
		except []int64
	}{
		{"int2[]int64", []int64{100}, []int64{100}},
		{"empty_string2[]int64", "[]", []int64{}},
		{"string2[]int64_case1", "[18,3,1,23]", []int64{18, 3, 1, 23}},
	}
	for _, tcase := range testCases {
		suite.Run(tcase.name, func() {
			value, err := Convert(tcase.value, ArrayInt64)
			assert.NoError(suite.T(), err)
			assert.ElementsMatch(suite.T(), tcase.except, value.([]int64))
		})
	}
}

func (suite *ConvertTestSuite) TestToArrayString() {
	testCases := []struct {
		name   string
		value  interface{}
		except []string
	}{
		{"empty_string2[]string", "[]", []string{}},
		{"string2[]int64", `["aa","bb"]`, []string{"aa", "bb"}},
	}
	for _, tcase := range testCases {
		suite.Run(tcase.name, func() {
			value, err := Convert(tcase.value, ArrayString)
			assert.NoError(suite.T(), err)
			assert.ElementsMatch(suite.T(), tcase.except, value.([]string))
		})
	}
}

func TestConvert(t *testing.T) {
	suite.Run(t, new(ConvertTestSuite))
}

// todo [optimize] other test
