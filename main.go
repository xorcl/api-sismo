package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xorcl/api-sismo/sismo"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
		if len(c.Errors) != 0 {
			// 2nd Try
			c.Next()
		}
	}
}

type Parser interface {
	GetRoute() string
	Parse(c *gin.Context)
}

func main() {
	parsers := []Parser{
		&sismo.RecentParser{},
		&sismo.HistoricParser{},
	}
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.Use(CORSMiddleware())
	for _, parser := range parsers {
		r.GET(fmt.Sprintf("/%s", parser.GetRoute()), parser.Parse)
	}
	r.Run()
}
