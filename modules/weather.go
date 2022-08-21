package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

var ApiKey string

type apidata = map[string]interface{}

type geodata struct {
	name string  `json:"name"`
	lat  float32 `json:"lat"`
	lon  float32 `json:"lon"`
}

type result struct {
	temp int32
	err  bool
}

type WeatherModule struct {
	city        string
	countryCode string
	err         bool
	ch          chan result
	temp        int32
}

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

func gettemp(ch chan result, lat, lon float32) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s",
		lat,
		lon,
		ApiKey,
	)

	data := make(apidata)
	if get(url, &data) {
		ch <- result{0, true}
	} else {
		ch <- result{
			int32(math.Round(data["main"].(apidata)["temp"].(float64))),
			true,
		}
	}
}

func fetch(ch chan result, lat, lon float32) {
	ticker := time.NewTicker(time.Minute)
	gettemp(ch, lat, lon)
	for _ = range ticker.C {
		gettemp(ch, lat, lon)
	}
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
	ch := make(chan result)
	weather := WeatherModule{
		city,
		countryCode,
		err,
		ch,
		0,
	}
	if err == false {
		go fetch(ch, location[0].lat, location[0].lon)
	}
	return &weather
}

func (c *WeatherModule) Output() string {
	if c.err == true {
		return BadOutput("missing geolocation")
	}
	select {
	case res := <-c.ch:
		c.temp = res.temp
	default:
	}
	return SimpleOutput(fmt.Sprintf("%d°C", c.temp))
}
