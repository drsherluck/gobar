package main

import (
	"fmt"
	"time"
	"./modules"
)

type Bar struct {
	modules []modules.Module
}

func NewBar() *Bar {
	clock := modules.CreateClock()
	network := modules.CreateNetwork("wlp5s0")
	battery := modules.CreateBattery()
	modules := []modules.Module{ network, battery, clock }
	return &Bar{ modules }
}

func (b *Bar) GenerateOutput() string {
	var output string = "["
	for i, m := range b.modules {
		if i == 0 {
			output = fmt.Sprintf("%s%s", output, m.Output())
			continue
		}
		output = fmt.Sprintf("%s,%s", output, m.Output())
	}
	return fmt.Sprintf("%s] ,", output)
}

func main() {
	fmt.Print("{\"version\":1}")
	fmt.Print("[")

	// Initialize Bar
	bar := NewBar()

	// Print loop
	for {
		fmt.Print(bar.GenerateOutput())
		time.Sleep(time.Second)
	}
}


