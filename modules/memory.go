package modules

import (
	"fmt"
	proc "github.com/c9s/goprocinfo/linux"
)

type MemoryModule struct{}

func Memory() *MemoryModule {
	return &MemoryModule{}
}

func (m *MemoryModule) Output() string {
	info, err := proc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return BadOutput("Error: MemoryInfo")
	}
	out := fmt.Sprintf("Memory %d", info.MemAvailable/1024)
	return SimpleOutput(out)
}
