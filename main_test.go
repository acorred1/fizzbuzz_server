package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFizzBuzzHTTP(t *testing.T) {
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
		{"str1=bar&str2=foo&limit=100&int1=-1&int2=3", "Error:Field validation for 'Int1' failed on the 'min' tag"},
		{"str1=bar&str2=foo&limit=100&int1=7&int2=-3", "Error:Field validation for 'Int2' failed on the 'min' tag"},

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

func TestGetStatus(t *testing.T) {
	// Mock prometheus response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, PrometheusQueryPath, r.URL.Path)
		assert.Equal(t, "topk(1,sum(http_requests_total)by(int1,int2,limit,str1,str2))", r.URL.Query().Get("query"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
{
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
		 }
		`))
	}))
	defer server.Close()
	PrometheusURL = server.URL

	router := setUpRouter()
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/stats", nil)
	router.ServeHTTP(resp, req)
	expectedBody := `{"top_fizzbuzz_request":{"request_params":{"int1":2,"int2":30,"limit":100,"str1":"fizz","str2":"oo"},"num_hits":5}}`

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, resp.Header().Get("Content-Type"), "application/json; charset=utf-8")
	assert.Equal(t, expectedBody, resp.Body.String())
}

func TestGetMetrics(t *testing.T) {
	router := setUpRouter()
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, resp.Header().Get("Content-Type"), "text/plain; version=0.0.4; charset=utf-8")
	assert.Contains(t, resp.Body.String(), PrometheusTopHTTPRequestMetricName)
}
