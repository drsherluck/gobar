package main

import (
	"fmt"
	"github.com/drsherluck/gobar/modules"
	"io"
	"os"
	"time"
)

func isEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return true, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

type Bar struct {
	modules []modules.Module
}

func NewBar() *Bar {
	modules := make([]modules.Module, 0, 10)
	return &Bar{modules}
}

// Adds a module to the module list
func (b *Bar) AddModule(module modules.Module) {
	b.modules = append(b.modules, module)
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

	return fmt.Sprintf("%s],", output)
}

func (b *Bar) Init() {
	fmt.Print("{\"version\":1}\n")
	fmt.Print("[")

	// Add modules in desired order
	// Last is the rightmost in the bar
	b.AddModule(modules.Network("enp4s0"))
	b.AddModule(modules.Volume())
	b.AddModule(modules.Memory())
	b.AddModule(modules.Weather())
	if ok, _ := isEmpty("/sys/class/power_supply"); ok == false {
		b.AddModule(modules.Battery())
	}
	b.AddModule(modules.Clock())
}

// Print loop
func (b *Bar) Run() {
	for {
		fmt.Print(b.GenerateOutput())
		time.Sleep(time.Second)
	}
}

// Entry
func main() {
	bar := NewBar()
	bar.Init()
	bar.Run()
}
