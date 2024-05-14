package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"
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
	last result
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

func gettemp(ch chan result, lat, lon float32, token string) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s",
		lat,
		lon,
		token,
	)
	data := make(apidata)
	if get(url, &data) {
		ch <- result{0, true}
	} else {
		t := math.Round(data["main"].(apidata)["temp"].(float64))
		ch <- result{int32(t), false}
	}
}

func fetch(ch chan result, lat, lon float32) {
	token := os.Getenv("OPEN_WEATHERMAP_API_KEY")
	ticker := time.NewTicker(time.Minute)
	gettemp(ch, lat, lon, token)
	for _ = range ticker.C {
		gettemp(ch, lat, lon, token)
	}
}

func Weather() *WeatherModule {
	url := fmt.Sprintf("http://ip-api.com/json/")
	loc := geodata{}
	err := get(url, &loc)

	// channel and ticker to fetch weather data
	ch := make(chan result)
	if err == false {
		go fetch(ch, loc.Lat, loc.Lon)
	}
	weather := WeatherModule{err, ch, result{0, false}}
	return &weather
}

func (c *WeatherModule) Output() string {
	if c.err {
		return BadOutput("missing geolocation")
	}
	select {
	case res := <-c.ch:
		c.last = res
	default:
	}
	if c.last.err {
		return BadOutput("error")
	}
	return SimpleOutput(fmt.Sprintf("%d°C", c.last.temp))
}
