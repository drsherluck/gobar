package modules 
import (
	"fmt"
	"github.com/vishvananda/netlink"
)

type NetworkModule struct {
	dev string
	
	// bytes received and transmitted
	rx uint64
	tx uint64

}

func Network(name string) *NetworkModule {
	module := NetworkModule{name, 0, 0}
	return &module
}

func (n *NetworkModule) activity(link netlink.Link) (uint64, uint64) {
	rx := link.Attrs().Statistics.RxBytes
	tx := link.Attrs().Statistics.TxBytes

	defer func(){
		n.rx = rx
		n.tx = tx
	}()

	return rx - n.rx, tx - n.tx
}

func readable(bytes uint64) string {
	if bytes > 1E6 {
		return fmt.Sprintf("%.1fMB", float64(bytes / 1E6))
	}
	return fmt.Sprintf("%.1fKB", float64(bytes / 1000))
}

func status(r, t uint64) string {
	return fmt.Sprintf("[%s, %s]", readable(r), readable(t))
}

func (n *NetworkModule) Output() string {
	link, err := netlink.LinkByName(n.dev)
	if err != nil {
		return BadOutput("disconnected")
	}
	
	activity := status(n.activity(link))
	return GoodOutput(activity)
}

	
