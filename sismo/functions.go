package sismo

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ParseEvents(response *Response, url string, c *gin.Context) {
	magnitude, ok := c.Request.URL.Query()["magnitude"]
	if !ok || len(magnitude) == 0 {
		magnitude = []string{"0"}
	}
	intMagnitude, err := strconv.Atoi(magnitude[0])
	if err != nil {
		response.SetStatus(12)
		logrus.WithFields(logrus.Fields{
			"error": response.StatusDescription,
		}).Error(err)
		c.JSON(500, &response)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		response.SetStatus(21)
		logrus.WithFields(logrus.Fields{
			"error": response.StatusDescription,
		}).Error(err)
		c.JSON(500, &response)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		response.SetStatus(20)
		logrus.WithFields(logrus.Fields{
			"error": response.StatusDescription,
		}).Error(err)
		c.JSON(500, &response)
		return
	}
	events := make([]Event, 0)
	doc.Find(TABLE_SELECTOR).Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return // skip header
		}
		event := Event{}
		s.Find("td").Each(func(j int, s *goquery.Selection) {
			switch j {
			case 0: // Local date and URL
				link, ok := s.Find("a").Attr("href")
				if !ok {
					logrus.Error("error parsing event url: %s", err)
					return
				}
				urlComponents := strings.Split(link, "/")
				event.URL = BASE_URL + link
				event.ID = strings.Split(urlComponents[len(urlComponents)-1], ".")[0]
				event.MapURL = fmt.Sprintf(
					"%s%s/map_img/%s.jpeg",
					BASE_URL,
					strings.Join(urlComponents[:len(urlComponents)-1], "/"),
					event.ID,
				)
				date := strings.TrimSpace(s.Text())
				event.LocalDate = strings.TrimSpace(date)
			case 1: // UTC Date
				date := strings.TrimSpace(s.Text())
				event.UTCDate = strings.TrimSpace(date)
			case 2: // Latitude
				latitude := strings.TrimSpace(s.Text())
				floatLatitude, err := strconv.ParseFloat(latitude, 64)
				if err != nil {
					logrus.Error("error parsing latitude: %s", err)
					return
				}
				event.Latitude = floatLatitude
			case 3: // Longitude
				longitude := strings.TrimSpace(s.Text())
				floatLongitude, err := strconv.ParseFloat(longitude, 64)
				if err != nil {
					logrus.Error("error parsing longitude: %s", err)
					return
				}
				event.Longitude = floatLongitude
			case 4: // Depth
				depth := strings.TrimSpace(s.Text())
				floatDepth, err := strconv.ParseFloat(depth, 64)
				if err != nil {
					logrus.Error("error parsing depth: %s", err)
					return
				}
				event.Depth = floatDepth
			case 5: // Magnitude
				magnitude := strings.TrimSpace(s.Text())
				magnitudeArr := strings.Split(magnitude, " ")
				if len(magnitudeArr) != 2 {
					logrus.Error("error parsing magnitude: malformed field %s", magnitude)
					return
				}
				floatMagnitude, err := strconv.ParseFloat(magnitudeArr[0], 64)
				if err != nil {
					logrus.Error("error parsing magnitude: %s", err)
					return
				}
				event.Magnitude = &Magnitude{
					Value:       floatMagnitude,
					MeasureUnit: magnitudeArr[1],
				}
			case 6: // Geographic Reference
				georeference := strings.TrimSpace(s.Text())
				event.GeoReference = strings.TrimSpace(georeference)
			}
		})
		events = append(events, event)
	})
	if intMagnitude != 0 {
		filteredEvents := make([]Event, 0)
		for _, event := range events {
			if event.Magnitude.Value >= float64(intMagnitude) {
				filteredEvents = append(filteredEvents, event)
			}
		}
		events = filteredEvents
	}
	response.Events = events
	response.SetStatus(0)
	c.JSON(200, response)
}
