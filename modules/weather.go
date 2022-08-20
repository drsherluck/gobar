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

type weather struct {
	data map[string]float32 `json:"main"`
}

type WeatherModule struct {
	city        string
	countryCode string
	lat         float32
	lon         float32
	err         bool
}

var ApiKey string
var oldtime = time.Now().Add(-time.Hour)
var temp = int(0)

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
	jerr := json.Unmarshal(body, &data)
	return jerr != nil
}

func getGeolocation(city, countryCode string) (float32, float32, bool) {
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s,%s&limit=1&appid=%s", city, countryCode, ApiKey)
	location := make([]geodata, 1)
	err := get(url, &location)
	return location[0].Lat, location[0].Lon, err
}

func Weather(city, countryCode string) *WeatherModule {
	lat, lon, err := getGeolocation(city, countryCode)
	weather := WeatherModule{city, countryCode, lat, lon, err}
	return &weather
}

func (c *WeatherModule) Output() string {
	if time.Now().Sub(oldtime).Minutes() < 1.0 {
		return SimpleOutput(fmt.Sprintf("%d°C", temp))
	}
	oldtime = time.Now()
	if c.err == true {
		lat, lon, err := getGeolocation(c.city, c.countryCode)
		c.lat = lat
		c.lon = lon
		c.err = err
		return BadOutput("missing geo")
	}
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s", c.lat, c.lon, ApiKey)

	data := make(map[string]interface{})
	if get(url, &data) {
		return BadOutput("error")
	}
	temp = int(math.Round(data["main"].(map[string]interface{})["temp"].(float64)))
	return SimpleOutput(fmt.Sprintf("%d°C", temp))
}
