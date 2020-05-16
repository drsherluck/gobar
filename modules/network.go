package modules 
import (
	"strings"
	"fmt"
	"log"
	"os/exec"
)

type NetworkModule struct {
	Interface string
	Name string
}

func CreateNetwork(name string) *NetworkModule {
	module := NetworkModule{name, "network"}
	return &module
}

func (n *NetworkModule) Output() string {
	out, err := exec.Command("nmcli", "-f", 
			"capabilities.carrier-detect,capabilities.speed", 
			"d", "show", n.Interface).Output()
	if err != nil {
		log.Fatal(err)
	}
	
	arr := strings.Fields(string(out))
	if arr[3] != "unknown"  {
		return GoodOutput(fmt.Sprintf("[%s Mbit/s]", arr[3]))
	}

	return BadOutput("disconnected")
}

	
