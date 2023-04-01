package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	PrometheusQueryPath                = "/api/v1/query"
	PrometheusTopHTTPRequestMetricName = "http_requests_total"
)

var PrometheusURL = "http://prometheus:9090"
var PrometheusTopHTTPRequestMetricLabels = []string{"int1", "int2", "limit", "str1", "str2"}

var (
	httpReqsMetric = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: PrometheusTopHTTPRequestMetricName,
			Help: "How many HTTP requests processed, partitioned by query params.",
		},
		PrometheusTopHTTPRequestMetricLabels,
	)
)

type FizzBuzzConfig struct {
	Int1  int    `form:"int1" json:"int1" binding:"required,numeric,max=2147483647,min=1"`
	Int2  int    `form:"int2" json:"int2" binding:"required,numeric,max=2147483647,min=1"`
	Limit int    `form:"limit" json:"limit" binding:"required,numeric,max=2147483647,min=1"`
	Str1  string `form:"str1" json:"str1" binding:"required"`
	Str2  string `form:"str2" json:"str2" binding:"required"`
}

type FizzBuzzResult struct {
	Result []string `json:"result"`
}

type StatsResponse struct {
	TopReq *FizzBuzzReqStats `json:"top_fizzbuzz_request"`
}

type FizzBuzzReqStats struct {
	Params  FizzBuzzConfig `json:"request_params"`
	NumHits int            `json:"num_hits"`
}

type PrometheusTopKMetricResponse struct {
	Data struct {
		Result []struct {
			Metric map[string]string `json:"metric"`
			Value  []string          `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
