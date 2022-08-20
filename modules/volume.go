package modules

import (
	"fmt"
	. "github.com/itchyny/volume-go"
)

type VolumeModule struct{}

func Volume() *VolumeModule {
	return &VolumeModule{}
}

func (m *VolumeModule) Output() string {
	muted, err := GetMuted()
	if err != nil {
		return BadOutput("?")
	}

	volume, err := GetVolume()
	if err != nil {
		return BadOutput("?")
	}

	if muted {
		return SimpleOutput("Muted")
	}
	return SimpleOutput(fmt.Sprintf("Vol %d", volume))
}
