package modules

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
	"net/http"
	"time"
)

type NetworkModule struct {
	dev string

	// bytes received and transmitted
	rx uint64
	tx uint64
	// connectivity check
	ch        chan bool
	connected bool
}

func connectivity(ch chan bool) {
	ticker := time.NewTicker(5 * time.Second)
	timeout := time.Duration(time.Second * 2)
	client := http.Client{
		Timeout: timeout,
	}
	for _ = range ticker.C {
		_, err := client.Get("http://google.com")
		ch <- err == nil
	}
}

func readable(bytes uint64) string {
	if bytes > 1e6 {
		return fmt.Sprintf("%.1fMB", float64(bytes/1e6))
	}
	return fmt.Sprintf("%.1fKB", float64(bytes/1000))
}

func status(r, t uint64) string {
	return fmt.Sprintf("[%s, %s]", readable(r), readable(t))
}

func ip(dev string) string {
	iface, _ := net.InterfaceByName(dev)
	addrs, err := iface.Addrs()
	ip := ""
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

func Network(name string) *NetworkModule {
	connected := make(chan bool)
	module := NetworkModule{name, 0, 0, connected, false}
	go connectivity(connected)
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

func (n *NetworkModule) Output() string {
	link, err := netlink.LinkByName(n.dev)
	if err != nil {
		return BadOutput("disconnected")
	}
	activity := fmt.Sprintf("%s %s", ip(n.dev), status(n.activity(link)))
	select {
	case status := <-n.ch:
		n.connected = status
	default:
	}
	if n.connected == false {
		return BadOutput(activity)
	}
	return GoodOutput(activity)
}
