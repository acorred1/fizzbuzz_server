package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func fizzbuzz(config FizzBuzzConfig) FizzBuzzResult {
	result := make([]string, 0)

	log.WithFields(
		log.Fields{"str1": config.Str1, "str2": config.Str2, "limit": config.Limit, "int1": config.Int1, "int2": config.Int2},
	).Debug("Computing FizzBuzz")
	for i := 1; i <= config.Limit; i++ {
		if i%config.Int1 == 0 && i%config.Int2 == 0 {
			result = append(result, fmt.Sprintf("%s%s", config.Str1, config.Str2))
		} else if i%config.Int1 == 0 {
			result = append(result, config.Str1)
		} else if i%config.Int2 == 0 {
			result = append(result, config.Str2)
		} else {
			result = append(result, fmt.Sprint(i))
		}
	}

	return FizzBuzzResult{Result: result}
}

func parsePrometheusResponse(body []byte) (int, *FizzBuzzConfig, error) {
	promTopKResp := PrometheusTopKMetricResponse{}
	json.Unmarshal([]byte(body), &promTopKResp)

	promResult := promTopKResp.Data.Result
	if len(promResult) == 0 {
		// No requests so far
		return 0, nil, nil
	} else if len(promResult) != 1 {
		errStr := fmt.Sprintf("Unexpected prometheus topk metrics result=%s", promResult)
		return 0, nil, errors.New(errStr)
	}

	numHitsStr := promResult[0].Value[1]
	numHits, err := strconv.Atoi(numHitsStr)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "Non integer type for request number of hits")
	}

	metricLabels := promResult[0].Metric

	int1, err := strconv.Atoi(metricLabels["int1"])
	if err != nil {
		return 0, nil, errors.Wrapf(err, "Non integer type for request int1 label")
	}

	int2, err := strconv.Atoi(metricLabels["int2"])
	if err != nil {
		return 0, nil, errors.Wrapf(err, "Non integer type for request int2 label")
	}
	limit, err := strconv.Atoi(metricLabels["limit"])
	if err != nil {
		return 0, nil, errors.Wrapf(err, "Non integer type for request limit label")
	}

	config := FizzBuzzConfig{Int1: int1, Int2: int2, Limit: limit, Str1: metricLabels["str1"], Str2: metricLabels["str2"]}

	return numHits, &config, nil
}
