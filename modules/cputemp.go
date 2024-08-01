package modules

import (
	"fmt"
	"math"
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
	files, err := os.ReadDir("/sys/class/hwmon")
	if err != nil {
		return &CpuTempModule{true, nil, nil, 0}
	}

	var dir string
	found_cpu := false
	for _, file := range files {
		hw_dir := fmt.Sprintf("/sys/class/hwmon/%s", file.Name())
		hw_files, err := os.ReadDir(hw_dir)
		if err != nil {
			break
		}

		for _, hf := range hw_files {
			if hf.Name() != "name" {
				continue
			}
			data, err := os.ReadFile(fmt.Sprintf("%s/name", hw_dir))
			if err != nil {
				break
			}
			name := strings.TrimSuffix(string(data), "\n")
			if name == "coretemp" || name == "k10temp" {
				found_cpu = true
				dir = hw_dir
				break
			}

		}
	}

	var zones []string
	if found_cpu {
		files, err = os.ReadDir(dir)
		if err != nil {
			return &CpuTempModule{true, nil, nil, 0}
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), "temp") && strings.HasSuffix(file.Name(), "_input") {
				zones = append(zones, fmt.Sprintf("%s/%s", dir, file.Name()))
			}
		}
	} else {
		files, err = os.ReadDir("/sys/class/thermal")
		if err != nil {
			return &CpuTempModule{true, nil, nil, 0}
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), "thermal_zone") && file.Name() != "thermal_zone0" {
				zones = append(zones, fmt.Sprintf("/sys/class/thermal/%s/temp", file.Name()))
			}
		}
	}

	ch := make(chan cpu_result)
	go fetchCpuTemp(ch, zones)
	return &CpuTempModule{false, ch, zones, 0.0}
}

func readtemp(zone string) (float64, error) {
	data, err := os.ReadFile(zone)
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
	var max = 0.0
	for _, zone := range zones {
		temp, err := readtemp(zone)
		if err != nil {
			ch <- cpu_result{0, true}
			return
		}
		max = math.Max(temp, max)
	}
	ch <- cpu_result{max, false}
}

func fetchCpuTemp(ch chan cpu_result, zones []string) {
	ticker := time.NewTicker(time.Second)
	getCpuTemp(ch, zones)
	for range ticker.C {
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
