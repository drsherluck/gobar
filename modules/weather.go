package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type apidata = map[string]interface{}

type geodata struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type result struct {
	temp int32
	err  error
}

func Err(err error) result {
	return result{0, err}
}

func Ok(temp int32) result {
	return result{temp, nil}
}

func (r *result) isOk() bool {
	return r.err == nil
}

func (r *result) isErr() bool {
	return r.err != nil
}

type WeatherModule struct {
	err  bool
	ch   chan result
	last result
}

func get(url string, data any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &data)
}

func get_temperature_from_response(data apidata) result {
	val, ok := data["error"].(bool)
	if ok && val {
		return Err(errors.New(data["reason"].(string)))
	}
	t := math.Round(data["current"].(apidata)["temperature_2m"].(float64))
	return Ok(int32(t))
}

func get_temperature(ch chan result, lat, lon float32) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
		lat,
		lon,
	)
	data := make(apidata)
	err := get(url, &data)
	if err != nil {
		ch <- Err(err)
	} else {
		ch <- get_temperature_from_response(data)
	}
}

func fetch(ch chan result, lat, lon float32) {
	ticker := time.NewTicker(time.Minute)
	get_temperature(ch, lat, lon)
	for _ = range ticker.C {
		get_temperature(ch, lat, lon)
	}
}

func Weather() *WeatherModule {
	url := fmt.Sprintf("http://ip-api.com/json/")
	loc := geodata{}
	err := get(url, &loc)

	// channel and ticker to fetch weather data
	ch := make(chan result)
	if err == nil {
		go fetch(ch, loc.Lat, loc.Lon)
	}
	weather := WeatherModule{err != nil, ch, Ok(0)}
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
	if c.last.isErr() {
		return BadOutput("error")
	}
	return SimpleOutput(fmt.Sprintf("%d°C", c.last.temp))
}
