package basic04_pinning_maps

import (
	"fmt"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

const (
	XdpPass   = uint32(2)
	bpfFSPath = "/sys/fs/bpf"
)

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

	// Name of the kernel function to trace.
	fn := "xdp_stats_map"

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	pinPath := path.Join(bpfFSPath, ifaceName, fn)
	if err := os.MkdirAll(pinPath, os.ModePerm); err != nil {
		log.Fatalf("failed to create bpf fs subpath: %+v", err)
	}

	// Load pre-compiled programs into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			// Pin the map to the BPF filesystem and configure the
			// library to automatically re-write it in the BPF
			// program so it can be re-used if it already exists or
			// create it if not
			PinPath: pinPath,
		},
	}); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpPassFunc,
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

	pmap, err := ebpf.LoadPinnedMap(pinPath, nil)
	if err != nil {
		log.Fatalf("could not load pinned map %s: %s", pinPath, err)
	}

	for range ticker.C {
		s, err := formatMapContents(pmap)
		if err != nil {
			log.Printf("Error reading map: %s", err)
			continue
		}
		log.Printf("Map contents:\n%s", s)
	}
}

func formatMapContents(m *ebpf.Map) (string, error) {
	var (
		sb   strings.Builder
		data []datarec
	)
	err := m.Lookup(XdpPass, &data)

	var sumPkts, sumBytes uint64
	for i := 0; i < len(data); i++ {
		sumPkts += data[i].RxPackets
		sumBytes += data[i].RxBytes
	}
	sb.WriteString(fmt.Sprintf("\tpackets => %d\n", sumPkts))
	sb.WriteString(fmt.Sprintf("\tbytes => %d\n", sumBytes))
	return sb.String(), err
}
