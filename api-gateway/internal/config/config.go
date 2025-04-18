package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port              string
	AuthServiceUrl    string
	MicroServiceUrl   string
	RedisServerAdress string
	RedisDBindex      int
	MaxRequests       int
	ConcurrencyLimit  int
	RateLimitTTL      time.Duration
}

func Load() *Config {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8081"
	}

	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://auth:8080"
	}

	microURL := os.Getenv("MICRO_SERVICE_URL")
	if microURL == "" {
		microURL = "http://functionservice:8082"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redisindex := 0
	redisindexString := os.Getenv("REDIS_DB")
	if redisindexString != "" {
		redisindexInt, err := strconv.Atoi(redisindexString)
		if err == nil {
			redisindex = redisindexInt
		}

	}

	maxrequests := 10
	ratelimitString := os.Getenv("RATE_LIMIT_REQ")
	if ratelimitString != "" {
		ratelimitInt, err := strconv.Atoi(ratelimitString)
		if err == nil {
			maxrequests = ratelimitInt
		}
	}

	rateTTL := time.Minute
	if v := os.Getenv("RATE_LIMIT_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			rateTTL = d
		}
	}

	concur := 5
	if v := os.Getenv("CONCURRENCY_LIMIT"); v != "" {
		if c, err := strconv.Atoi(v); err == nil {
			concur = c
		}
	}

	return &Config{
		Port:              port,
		AuthServiceUrl:    authURL,
		MicroServiceUrl:   microURL,
		RedisServerAdress: redisAddr,
		RedisDBindex:      redisindex,
		MaxRequests:       maxrequests,
		RateLimitTTL:      rateTTL,
		ConcurrencyLimit:  concur,
	}
}
