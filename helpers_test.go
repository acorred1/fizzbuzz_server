package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFizzBuzz(t *testing.T) {
	cases := []struct {
		config   FizzBuzzConfig
		expected []string
	}{
		{
			FizzBuzzConfig{Int1: 2, Int2: 3, Limit: 10, Str1: "fizz", Str2: "buzz"},
			[]string{"1", "fizz", "buzz", "fizz", "5", "fizzbuzz", "7", "fizz", "buzz", "fizz"},
		},
		{
			FizzBuzzConfig{Int1: 2, Int2: 3, Limit: 0, Str1: "fizz", Str2: "buzz"},
			[]string{},
		},
		{
			FizzBuzzConfig{Int1: -2, Int2: 3, Limit: 6, Str1: "abc", Str2: "123"},
			[]string{"1", "abc", "123", "abc", "5", "abc123"},
		},
		{
			FizzBuzzConfig{Int1: 2, Int2: 2, Limit: 10, Str1: "foo", Str2: "bar"},
			[]string{"1", "foobar", "3", "foobar", "5", "foobar", "7", "foobar", "9", "foobar"},
		},
	}

	for _, testCase := range cases {
		actual := fizzbuzz(testCase.config).Result
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestParsePrometheusResponse(t *testing.T) {
	cases := []struct {
		respBody              []byte
		expectedNumHits       int
		expectedFizzBuzzCofig *FizzBuzzConfig
	}{
		{
			[]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`),
			0,
			nil,
		},
		{
			[]byte(
				`{
				"status":"success",
				"data":{
					 "resultType":"vector",
					 "result":[
					 {
					 "metric":{"int1":"2","int2":"30","limit":"100","str1":"fizz","str2":"oo"},
					 "value":[1680208159.058,"5"]
				 }
				 ]
			 }
		 }`,
			),
			5,
			&FizzBuzzConfig{Int1: 2, Int2: 30, Limit: 100, Str1: "fizz", Str2: "oo"},
		},
	}

	for _, testCase := range cases {
		actualNumHits, actualFizzBuzzConfig, actualErr := parsePrometheusResponse(testCase.respBody)
		if assert.NoError(t, actualErr) {
			assert.Equal(t, testCase.expectedNumHits, actualNumHits)
			assert.Equal(t, testCase.expectedFizzBuzzCofig, actualFizzBuzzConfig)
		}
	}
}

func TestParsePrometheusResponse_Error(t *testing.T) {
	cases := []struct {
		respBody    []byte
		expectedErr string
	}{
		{
			[]byte(
				`
				{
					"status":"success",
					"data":{
						 "resultType":"vector",
						 "result":[{},{}]
				 }
			 }
			 `,
			),
			"Unexpected prometheus topk metrics result=[{map[] []} {map[] []}]",
		},
		{
			[]byte(
				`{
				"status":"success",
				"data":{
					 "resultType":"vector",
					 "result":[
					 {
					 "value":[1680208159.058,"abc"]
				 }
				 ]
			 }
		 }`,
			),
			"Non integer type for request number of hits: strconv.Atoi: parsing \"abc\": invalid syntax",
		},
		{
			[]byte(
				`{
				"status":"success",
				"data":{
					 "resultType":"vector",
					 "result":[
					 {
					 "metric":{"int1":"A","int2":"30","limit":"100","str1":"fizz","str2":"oo"},
					 "value":[1680208159.058,"5"]
				 }
				 ]
			 }
		 }`,
			),
			"Non integer type for request int1 label: strconv.Atoi: parsing \"A\": invalid syntax",
		},
		{
			[]byte(
				`{
				"status":"success",
				"data":{
					 "resultType":"vector",
					 "result":[
					 {
					 "metric":{"int1":"2","int2":"30A","limit":"100","str1":"fizz","str2":"oo"},
					 "value":[1680208159.058,"5"]
				 }
				 ]
			 }
		 }`,
			),
			"Non integer type for request int2 label: strconv.Atoi: parsing \"30A\": invalid syntax",
		},
		{
			[]byte(
				`{
				"status":"success",
				"data":{
					 "resultType":"vector",
					 "result":[
					 {
					 "metric":{"int1":"2","int2":"30","limit":"abcd","str1":"fizz","str2":"oo"},
					 "value":[1680208159.058,"5"]
				 }
				 ]
			 }
		 }`,
			),
			"Non integer type for request limit label: strconv.Atoi: parsing \"abcd\": invalid syntax",
		},
	}

	for _, testCase := range cases {
		actualNumHits, actualFizzBuzzConfig, actualErr := parsePrometheusResponse(testCase.respBody)
		assert.Equal(t, 0, actualNumHits)
		assert.Nil(t, actualFizzBuzzConfig)
		assert.EqualError(t, actualErr, testCase.expectedErr)
	}
}
