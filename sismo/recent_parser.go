package sismo

import (
	"github.com/gin-gonic/gin"
)

const RECENT_URL = BASE_URL + "/ultimos_sismos.html"

type RecentParser struct{}

func (bp *RecentParser) GetRoute() string {
	return "recent"
}

func (bp *RecentParser) Parse(c *gin.Context) {
	response := &Response{}
	ParseEvents(response, RECENT_URL, c)
}
