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

func ParseEvents(response *Response, urls []string, c *gin.Context) {
	events := make([]Event, 0)
	magnitude, ok := c.Request.URL.Query()["magnitude"]
	if !ok || len(magnitude) == 0 {
		magnitude = []string{"0"}
	}
	limit, ok := c.Request.URL.Query()["limit"]
	if !ok || len(limit) == 0 {
		limit = []string{"0"}
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
	intLimit, err := strconv.Atoi(limit[0])
	if err != nil {
		response.SetStatus(12)
		logrus.WithFields(logrus.Fields{
			"error": response.StatusDescription,
		}).Error(err)
		c.JSON(500, &response)
		return
	}
	for _, url := range urls {

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
		doc.Find(TABLE_SELECTOR).Each(func(i int, s *goquery.Selection) {
			if i < 2 {
				return // skip header and search field
			}
			event := Event{}
			s.Find("td").Each(func(j int, s *goquery.Selection) {
				switch j {
				case 0: // Local date, URL and Geo Reference
					link, ok := s.Find("a").Attr("href")
					if !ok {
						logrus.Errorf("error parsing event url: %s", err)
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
					dateLocation := strings.TrimSpace(s.Text())
					event.LocalDate = strings.TrimSpace(dateLocation[:19])
					event.GeoReference = strings.TrimSpace(dateLocation[19:])
				case 1: // UTC Date
					date := strings.TrimSpace(s.Text())
					event.UTCDate = strings.TrimSpace(date)
				case 2: // Latitude and Longitude
					latLng := strings.Split(strings.TrimSpace(s.Text()), " ")
					if len(latLng) != 2 {
						logrus.Errorf("error parsing latitude and longitude: %s", s.Text())
						return
					}
					latitude := strings.TrimSpace(latLng[0])
					floatLatitude, err := strconv.ParseFloat(latitude, 64)
					if err != nil {
						logrus.Errorf("error parsing latitude: %s", err)
						return
					}
					event.Latitude = floatLatitude
					longitude := strings.TrimSpace(latLng[1])
					floatLongitude, err := strconv.ParseFloat(longitude, 64)
					if err != nil {
						logrus.Errorf("error parsing longitude: %s", err)
						return
					}
					event.Longitude = floatLongitude
				case 3: // Depth
					depthArr := strings.Split(strings.TrimSpace(s.Text()), " ")
					if len(depthArr) != 2 {
						logrus.Errorf("error parsing latitude and longitude: %s", s.Text())
						return
					}
					floatDepth, err := strconv.ParseFloat(depthArr[0], 64)
					if err != nil {
						logrus.Errorf("error parsing depth: %s", err)
						return
					}
					event.Depth = floatDepth
				case 4: // Magnitude
					magnitude := strings.TrimSpace(s.Text())
					magnitudeArr := strings.Split(magnitude, " ")
					if len(magnitudeArr) != 2 {
						logrus.Errorf("error parsing magnitude: malformed field %s", magnitude)
						return
					}
					floatMagnitude, err := strconv.ParseFloat(magnitudeArr[0], 64)
					if err != nil {
						logrus.Errorf("error parsing magnitude: %s", err)
						return
					}
					event.Magnitude = &Magnitude{
						Value:       floatMagnitude,
						MeasureUnit: magnitudeArr[1],
					}
				}
			})
			events = append(events, event)
		})
	}
	if intMagnitude != 0 {
		filteredEvents := make([]Event, 0)
		for _, event := range events {
			if event.Magnitude.Value >= float64(intMagnitude) {
				filteredEvents = append(filteredEvents, event)
			}
		}
		events = filteredEvents
	}
	if intLimit > 0 && intLimit < len(events) {
		events = events[:intLimit]
	}
	response.Events = events
	response.SetStatus(0)
	c.JSON(200, response)
}
