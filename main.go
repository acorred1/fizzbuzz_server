package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func setUpRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	// Turn panics into 500s
	router.Use(gin.Recovery())

	router.GET("/api/v1/fizzbuzz", fizzbuzzHandler)
	router.GET("/api/v1/stats", statsHandler)

	// Set up prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}

func setUpLogger() {
	// This could be JSONFormatter if the logs were to be consumed by some log aggregator
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	levelStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		levelStr = "INFO"
	}

	logLevel, err := log.ParseLevel(levelStr)
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}

func init() {
	setUpLogger()
}

func main() {
	router := setUpRouter()
	router.Run()
}
