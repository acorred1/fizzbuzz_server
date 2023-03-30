package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetFizzBuzzHTTP_Success(t *testing.T) {
	router := setUpRouter()
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fizzbuzz?str1=foo&str2=bar&limit=25&int1=7&int2=3", nil)
	router.ServeHTTP(resp, req)
	expectedBody := `{"result":["1","2","bar","4","5","bar","foo","8","bar","10","11","bar","13","foo","bar","16","17","bar","19","20","foobar","22","23","bar","25"]}`

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, resp.Header().Get("Content-Type"), "application/json; charset=utf-8")
	assert.Equal(t, expectedBody, resp.Body.String())
}

func TestQueryParamValidation(t *testing.T) {
	router := setUpRouter()

	cases := []struct {
		queryParams   string
		expectedError string
	}{
		// Required
		{"str2=bar&limit=25&int1=7&int2=3", "Error:Field validation for 'Str1' failed on the 'required' tag"},
		{"str1=bar&limit=25&int1=7&int2=3", "Error:Field validation for 'Str2' failed on the 'required' tag"},
		{"str1=bar&str2=foo&int1=7&int2=3", "Error:Field validation for 'Limit' failed on the 'required' tag"},
		{"str1=bar&str2=foo&limit=100&int2=3", "Error:Field validation for 'Int1' failed on the 'required' tag"},
		{"str1=bar&str2=foo&limit=100&int1=7", "Error:Field validation for 'Int2' failed on the 'required' tag"},

		// Max values
		{"str1=bar&str2=foo&limit=2147483657&int1=7&int2=3", "Error:Field validation for 'Limit' failed on the 'max' tag"},
		{"str1=bar&str2=foo&limit=100&int1=2147483657&int2=3", "Error:Field validation for 'Int1' failed on the 'max' tag"},
		{"str1=bar&str2=foo&limit=100&int1=7&int2=2147483657", "Error:Field validation for 'Int2' failed on the 'max' tag"},

		// Min value
		{"str1=bar&str2=foo&limit=-12&int1=7&int2=3", "Error:Field validation for 'Limit' failed on the 'min' tag"},
		{"str1=bar&str2=foo&limit=100&int1=-2147483657&int2=3", "Error:Field validation for 'Int1' failed on the 'min' tag"},
		{"str1=bar&str2=foo&limit=100&int1=7&int2=-2147483657", "Error:Field validation for 'Int2' failed on the 'min' tag"},

		// Numeric
		{"str1=bar&str2=foo&limit=ABC&int1=7&int2=3", "strconv.ParseInt: parsing \\\"ABC\\\": invalid syntax"},
		{"str1=bar&str2=foo&limit=100&int1=FOO&int2=3", "strconv.ParseInt: parsing \\\"FOO\\\": invalid syntax"},
		{"str1=bar&str2=foo&limit=100&int1=7&int2=BAR", "strconv.ParseInt: parsing \\\"BAR\\\": invalid syntax"},
	}

	for _, testCase := range cases {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", "/api/v1/fizzbuzz?", testCase.queryParams), nil)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, resp.Header().Get("Content-Type"), "application/json; charset=utf-8")
		assert.Contains(t, resp.Body.String(), testCase.expectedError)
	}
}
