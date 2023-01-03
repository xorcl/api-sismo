package sismo

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HistoricParser struct{}

func (bp *HistoricParser) GetRoute() string {
	return "historic/:date"
}

func (bp *HistoricParser) Parse(c *gin.Context) {
	response := &Response{}
	date := c.Param("date")
	if len(date) == 0 {
		response.SetStatus(11)
		logrus.WithFields(logrus.Fields{
			"error": response.StatusDescription,
		}).Error(fmt.Errorf(response.StatusDescription))
		c.JSON(400, &response)
		return
	}
	url, err := getURLByDate(date)
	if err != nil {
		response.SetStatus(11)
		logrus.WithFields(logrus.Fields{
			"error": response.StatusDescription,
		}).Error(fmt.Errorf(response.StatusDescription))
		c.JSON(400, &response)
		return
	}
	ParseEvents(response, []string{url}, c)
}

func getURLByDate(date string) (string, error) {
	timeObj, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(HISTORIC_URL, timeObj.Year(), timeObj.Month(), timeObj.Year(), timeObj.Month(), timeObj.Day()), nil
}
