package main

import (
  "fmt"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
)

// TODO error wrapping
// TODO logging
// TODO tests
// TODO make port configurable via env variable

const (
  LimitDefault = "100"
  Int1Default = "3"
  Int2Default = "5"
  Str1Default = "fizz"
  Str2Default = "buzz"
)

type FizzBuzzConfig struct {
  Int1, Int2, Limit int
  Str1, Str2 string
}

type FizzBuzzResult struct {
	Result []string `json:"result"`
}


func fizzbuzz(config FizzBuzzConfig) FizzBuzzResult {
  result := make([]string, 0)

  for i := 1; i <= config.Limit; i++ {
    if i % (config.Int1 * config.Int2) == 0 {
      result = append(result, config.Str1 + config.Str2)
    } else if i % config.Int1 == 0 {
      result = append(result, config.Str1)
    } else if i % config.Int2 == 0 {
      result = append(result, config.Str2)
    } else {
      result = append(result, fmt.Sprint(i))
    }
  }

  return FizzBuzzResult{Result: result}
}

func parseAndValidateQueryParams(ctx *gin.Context) (*FizzBuzzConfig, error) {
  limitStr := ctx.DefaultQuery("limit", LimitDefault)
  int1Str := ctx.DefaultQuery("int1", Int1Default)
  int2Str := ctx.DefaultQuery("int2", Int2Default)
  str1 := ctx.DefaultQuery("str1", Str1Default)
  str2 := ctx.DefaultQuery("str2", Str2Default)

  limit, err := strconv.Atoi(limitStr)
  if err != nil {
    return nil, err
  }

  int1, err := strconv.Atoi(int1Str)
  if err != nil {
    return nil, err
  }

  int2, err := strconv.Atoi(int2Str)
  if err != nil {
    return nil, err
  }

  config := FizzBuzzConfig{
    Limit: limit,
    Int1: int1,
    Int2: int2,
    Str1: str1,
    Str2: str2,
  }

  return &config, nil
}

func fizzbuzzHandler(ctx *gin.Context) {
  config, err := parseAndValidateQueryParams(ctx)
  if err != nil {
    // TODO return an appropriate HTTP error
  }

  ctx.JSON(http.StatusOK, fizzbuzz(*config))
}

func main() {
  router := gin.New()

  router.GET("/api/v1/fizzbuzz", fizzbuzzHandler)

  router.Run()
}
