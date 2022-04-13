package middleware

import (
	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/gin-gonic/gin"
)

var (
	config = &statsd.ClientConfig{
		Address: "127.0.0.1:8125",
		Prefix:  "req-client",
	}
	client, _ = statsd.NewClientWithConfig(config)
)

func Counter() gin.HandlerFunc {
	return func(c *gin.Context) {
		client.Inc("req", 1, 1.0)
		c.Next()
	}
}

func CloseClient() {
	client.Close()
}
