package modules

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type cpu_result struct {
	temp float64
	err  bool
}

type CpuTempModule struct {
	err   bool
	ch    chan cpu_result
	zones []string
	temp  float64
}

func CpuTemp() *CpuTempModule {
	files, err := ioutil.ReadDir("/sys/class/thermal")
	has_err := false
	if err != nil {
		has_err = true
	}

	var zones []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "thermal_zone") && file.Name() != "thermal_zone0" {
			zones = append(zones, file.Name())
		}
	}
	ch := make(chan cpu_result)
	if has_err == false {
		go fetchCpuTemp(ch, zones)
	}
	return &CpuTempModule{has_err, ch, zones, 0.0}
}

func readtemp(zone string) (float64, error) {
	data, err := os.ReadFile(fmt.Sprintf("/sys/class/thermal/%s/temp", zone))
	if err != nil {
		return 0.0, err
	}
	number_string := strings.TrimSuffix(string(data), "\n")
	val, err := strconv.ParseInt(number_string, 10, 0)
	if err != nil {
		return 0.0, err
	}
	return float64(val) / 1000.0, nil
}

func getCpuTemp(ch chan cpu_result, zones []string) {
	var avg = 0.0
	for _, zone := range zones {
		temp, err := readtemp(zone)
		if err != nil {
			ch <- cpu_result{0, true}
			return
		}
		avg += temp
	}
	ch <- cpu_result{avg / float64(len(zones)), false}
}

func fetchCpuTemp(ch chan cpu_result, zones []string) {
	ticker := time.NewTicker(time.Second * 6)
	getCpuTemp(ch, zones)
	for _ = range ticker.C {
		getCpuTemp(ch, zones)
	}
}

func (m *CpuTempModule) Output() string {
	if m.err == true {
		return BadOutput("CPU ?")
	}
	select {
	case res := <-m.ch:
		m.temp = res.temp
	default:
	}
	return SimpleOutput(fmt.Sprintf("CPU %.2f°C", m.temp))
}
