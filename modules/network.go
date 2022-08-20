package modules

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
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

	defer func() {
		n.rx = rx
		n.tx = tx
	}()

	return rx - n.rx, tx - n.tx
}

func readable(bytes uint64) string {
	if bytes > 1e6 {
		return fmt.Sprintf("%.1fMB", float64(bytes/1e6))
	}
	return fmt.Sprintf("%.1fKB", float64(bytes/1000))
}

func status(ip string, r, t uint64) string {
	return fmt.Sprintf("%s [%s, %s]", ip, readable(r), readable(t))
}

func (n *NetworkModule) Output() string {
	link, err := netlink.LinkByName(n.dev)
	if err != nil {
		return BadOutput("disconnected")
	}
	iface, _ := net.InterfaceByName(n.dev)
	addrs, err := iface.Addrs()
	ip := ""
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	r, t := n.activity(link)
	activity := status(ip, r, t)
	return GoodOutput(activity)
}
