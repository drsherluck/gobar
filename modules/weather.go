package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

var (
	OWM_API_KEY string
)

type apidata = map[string]interface{}

type geodata struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type result struct {
	temp int32
	err  bool
}

type WeatherModule struct {
	err  bool
	ch   chan result
	temp int32
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
		OWM_API_KEY,
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

func Weather() *WeatherModule {
	url := fmt.Sprintf("http://ip-api.com/json/")
	location := geodata{}
	err := get(url, &location)

	// channel and ticker to fetch weather data
	ch := make(chan result)
	if err == false {
		go fetch(ch, location.Lat, location.Lon)
	}
	weather := WeatherModule{err, ch, 0}
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
