package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type geodata struct {
	Name string  `json:"name"`
	Lat  float32 `json:"lat"`
	Lon  float32 `json:"lon"`
}

type WeatherModule struct {
	city        string
	countryCode string
	lat         float32
	lon         float32
	err         bool
}

var (
	ApiKey  string
	oldtime = time.Now().Add(-time.Hour)
	temp    = int(0)
)

func get(url string, data any) bool {
	resp, err := http.Get(url)
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return true
	}
	err = json.Unmarshal(body, &data)
	return err != nil
}

func Weather(city, countryCode string) *WeatherModule {
	url := fmt.Sprintf(
		"http://api.openweathermap.org/geo/1.0/direct?q=%s,%s&limit=1&appid=%s",
		city,
		countryCode,
		ApiKey,
	)
	location := make([]geodata, 1)
	err := get(url, &location)
	weather := WeatherModule{city, countryCode, location[0].Lat, location[0].Lon, err}
	return &weather
}

func (c *WeatherModule) Output() string {
	if time.Now().Sub(oldtime).Minutes() >= 1.0 {
		oldtime = time.Now()
		if c.err == true {
			return BadOutput("missing geolocation")
		}
		url := fmt.Sprintf(
			"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s",
			c.lat,
			c.lon,
			ApiKey,
		)

		data := make(map[string]interface{})
		if get(url, &data) {
			return BadOutput("error")
		}
		temp = int(math.Round(data["main"].(map[string]interface{})["temp"].(float64)))
	}
	return SimpleOutput(fmt.Sprintf("%d°C", temp))
}
