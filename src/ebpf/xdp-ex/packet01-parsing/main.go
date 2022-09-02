package packet01_parsing

import (
	"fmt"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const XdpPass = uint32(2)

type datarec struct {
	RxPackets uint64
	RxBytes   uint64
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please specify a network interface")
	}

	// Look up the network interface by name.
	ifaceName := os.Args[1]
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("lookup network iface %q: %s", ifaceName, err)
	}

	// Load pre-compiled programs into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpIpFilter,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("could not attach XDP program: %s", err)
	}
	defer l.Close()

	log.Printf("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)
	log.Printf("Press Ctrl-C to exit and remove the program")

	// Print the contents of the BPF hash map (source IP address -> packet count).
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s, err := formatMapContents(objs.XdpStatsMap)
		if err != nil {
			log.Printf("Error reading map: %s", err)
			continue
		}
		log.Printf("Map contents:\n%s", s)
	}
}

func formatMapContents(m *ebpf.Map) (string, error) {
	var (
		sb      strings.Builder
		data    []datarec
		pktSum  uint64
		byteSum uint64
	)
	err := m.Lookup(XdpPass, &data)

	for _, datum := range data {
		pktSum += datum.RxPackets
		byteSum += datum.RxBytes
	}
	sb.WriteString(fmt.Sprintf("\tpackets => %d\n", pktSum))
	sb.WriteString(fmt.Sprintf("\tbytes => %d\n", byteSum))
	return sb.String(), err
}
