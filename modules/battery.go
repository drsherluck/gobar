package modules

import (
	"fmt"
	"github.com/distatus/battery"
)

type BatteryModule struct {
	state   battery.State
	current float64
	full    float64
}

func Battery() *BatteryModule {
	module := BatteryModule{battery.Unknown, 0.0, 1.0}
	module.update()
	return &module
}

func (b *BatteryModule) update() {
	battery, err := battery.Get(0)
	if err != nil {
		return
	}
	b.state = battery.State
	b.current = battery.Current
	b.full = battery.Full
}

func (b *BatteryModule) charge() int {
	return int((b.current / b.full) * 100)
}

func (b *BatteryModule) Output() string {
	b.update()

	// Create output string
	charge := b.charge()
	out := fmt.Sprintf("%d%s", charge, "%")

	switch b.state {
	case battery.Charging:
		return GoodOutput(out)
	case battery.Unknown:
		return BadOutput("unkown")
	default:
		if charge < 10 {
			return BadOutput(out)
		}
		return SimpleOutput(out)

	}
}
