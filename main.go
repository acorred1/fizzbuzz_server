package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// TODO consider refactoring into models, handlers, fizzbuzz and move tests too

type FizzBuzzConfig struct {
	Int1  int    `form:"int1" binding:"required,numeric,max=2147483647,min=1"`
	Int2  int    `form:"int2" binding:"required,numeric,max=2147483647,min=1"`
	Limit int    `form:"limit" binding:"required,numeric,max=2147483647,min=1"`
	Str1  string `form:"str1" binding:"required"`
	Str2  string `form:"str2" binding:"required"`
}

type FizzBuzzResult struct {
	Result []string `json:"result"`
}

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

func fizzbuzzHandler(ctx *gin.Context) {
	config := FizzBuzzConfig{}
	if err := ctx.ShouldBind(&config); err != nil {
		log.Error("Error parsing query params: ", err)
		// TODO consider improving these error messages
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": strings.Split(err.Error(), "\n")})
		return
	}

	ctx.JSON(http.StatusOK, fizzbuzz(config))
}

func setUpRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	// Turn panics into 500s
	router.Use(gin.Recovery())

	router.GET("/api/v1/fizzbuzz", fizzbuzzHandler)

	return router
}

func setUpLogger() {
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
