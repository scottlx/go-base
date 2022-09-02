package basic02_prog_by_name

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Please specify a network interface and packet action")
	}

	action := os.Args[2]

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

	// Attach the program.
	var prog *ebpf.Program
	switch action {
	case "drop":
		prog = objs.XdpDropFunc
	case "pass":
		prog = objs.XdpPassFunc
	default:
		log.Printf("Invalid packect option")
		return
	}

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   prog,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("could not attach XDP program: %s", err)
	}
	defer l.Close()

	log.Printf("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)
	log.Printf("The program runs 30 seconds")
	log.Printf("Press Ctrl-C to exit and remove the program")

	time.Sleep(time.Second * 30)
}
