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
	modules := make([]modules.Module, 0, 10)
	return &Bar{ modules }
}

// Adds a module to the module list
func (b *Bar)  AddModule(module modules.Module) {
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

	return fmt.Sprintf("%s] ,", output)
}

func (b *Bar) Init() {
	fmt.Print("{\"version\":1}")
	fmt.Print("[")

	// Add modules in desired order
	// Last is the rightmost in the bar
	b.AddModule(modules.Network("enp4s0"))
	b.AddModule(modules.Volume())
	b.AddModule(modules.Memory())
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


