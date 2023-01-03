package sismo

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type RecentParser struct{}

func (bp *RecentParser) GetRoute() string {
	return "recent"
}

func (bp *RecentParser) Parse(c *gin.Context) {
	response := &Response{}
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	recentUrls := []string{
		fmt.Sprintf(HISTORIC_URL, today.Year(), today.Month(), today.Year(), today.Month(), today.Day()),
		fmt.Sprintf(HISTORIC_URL, yesterday.Year(), yesterday.Month(), yesterday.Year(), yesterday.Month(), yesterday.Day()),
	}
	ParseEvents(response, recentUrls, c)
}
