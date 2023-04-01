package main

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)


func fizzbuzzHandler(ctx *gin.Context) {
	config := FizzBuzzConfig{}
	if err := ctx.ShouldBind(&config); err != nil {
		log.Errorf("Error parsing query params: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": strings.Split(err.Error(), "\n")})
		return
	}

	httpReqsMetric.WithLabelValues(
		fmt.Sprint(config.Int1), fmt.Sprint(config.Int2), fmt.Sprint(config.Limit), config.Str1, config.Str2,
	).Inc()
	ctx.JSON(http.StatusOK, fizzbuzz(config))
}

func statsHandler(ctx *gin.Context) {
	// Query the top request from prometheus
	topRequestQuery := fmt.Sprintf(
		"topk(1,sum(%s)by(%s))",
		PrometheusTopHTTPRequestMetricName,
		strings.Join(PrometheusTopHTTPRequestMetricLabels, ","),
	)
	requestURL := fmt.Sprintf(
		"%s%s?query=%s",
		PrometheusURL,
		PrometheusQueryPath,
		topRequestQuery,
	)

	res, err := http.Get(requestURL)
	if err != nil {
		log.Errorf("Error requesting metric from prometheus: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Failed to read prometheus response body. Err=%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	numHits, fizzBuzzConfig, err := parsePrometheusResponse(body)
	if err != nil {
		log.Errorf("Failed to parse prometheus TopK metric. Err=%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var stats StatsResponse
	if numHits == 0 {
		stats = StatsResponse{TopReq: nil}
	} else {
		stats = StatsResponse{
			TopReq: &FizzBuzzReqStats{
				NumHits: numHits,
				Params: *fizzBuzzConfig,
			},
		}
	}

	ctx.JSON(http.StatusOK, stats)
}
