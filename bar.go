package main

import (
	"fmt"
	"github.com/drsherluck/gobar/modules"
	"github.com/mkideal/cli"
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

func (b *Bar) Init(cfg *Config) {
	for _, m := range cfg.Modules {
		switch m {
		case "network":
			b.AddModule(modules.Network(cfg.Network.Interface))
		case "volume":
			b.AddModule(modules.Volume())
		case "cputemp":
			b.AddModule(modules.CpuTemp())
		case "memory":
			b.AddModule(modules.Memory())
		case "weather":
			b.AddModule(modules.Weather())
		case "battery":
			b.AddModule(modules.Battery())
		case "time":
			b.AddModule(modules.Clock())
		}
	}
}

// Print loop
func (b *Bar) Run() {
	fmt.Print("{\"version\":1}\n")
	fmt.Print("[")
	for {
		fmt.Print(b.GenerateOutput())
		time.Sleep(time.Second)
	}
}

func setupcfg(argv *argT) *Config {
	var cfg *Config
	var err error
	if len(argv.Config) > 0 {
		if _, err := os.Stat(argv.Config); err != nil {
			panic(err)
		}
		cfg, err = NewCustomConfig(argv.Config)
		if err != nil {
			panic(err)
		}
	} else if dirname, err := os.UserHomeDir(); err == nil {
		path := fmt.Sprintf("%s/.config/gobar/config.toml", dirname)
		if _, err = os.Stat(path); err == nil {
			cfg, err = NewCustomConfig(path)
			if err != nil {
				panic(err)
			}
		}
	}

	if cfg == nil {
		cfg = NewDefaultConfig()
	}

	if len(argv.NIC) > 0 {
		cfg.Network.Interface = argv.NIC
	}
	if argv.ExcludeBattery {
		temp := make([]string, 0, len(cfg.Modules)-1)
		for _, s := range cfg.Modules {
			if s != "battery" {
				temp = append(temp, s)
			}
		}
		cfg.Modules = temp
	}
	return cfg
}

type argT struct {
	cli.Helper
	Config         string `cli:"config" usage:"the path to the config file"`
	NIC            string `cli:"nic" usage:"the network interface to poll from"`
	ExcludeBattery bool   `cli:"no-bat" usage:"excludes the battery module" dft:"false"`
}

// Entry
func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		cfg := setupcfg(argv)
		bar := NewBar()
		bar.Init(cfg)
		bar.Run()
		return nil
	}))
}
