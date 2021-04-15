package wol

import (
	"errors"
	"fmt"
	"net"
)

////////////////////////////////////////////////////////////////////////////////

const (
	dbPath = "/.config/go-wol/bolt.db"
)

var (
	// Define holders for the cli arguments we wish to parse.
	cliFlags struct {
		Version            bool   `short:"v" long:"version"`
		Help               bool   `short:"h" long:"help"`
		NoColor            bool   `short:"n" long:"no-color"`
		BroadcastInterface string `short:"i" long:"interface" default:""`
		BroadcastIP        string `short:"b" long:"bcast" default:"255.255.255.255"`
		UDPPort            string `short:"p" long:"port" default:"9"`
	}
)

////////////////////////////////////////////////////////////////////////////////

// ipFromInterface returns a `*net.UDPAddr` from a network interface name.
func ipFromInterface(iface string) (*net.UDPAddr, error) {
	ief, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	addrs, err := ief.Addrs()
	if err == nil && len(addrs) <= 0 {
		err = fmt.Errorf("no address associated with interface %s", iface)
	}
	if err != nil {
		return nil, err
	}

	// Validate that one of the addrs is a valid network IP address.
	for _, addr := range addrs {
		switch ip := addr.(type) {
		case *net.IPNet:
			if !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				return &net.UDPAddr{
					IP: ip.IP,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no address associated with interface %s", iface)
}

// Run the wake command.
func WakeCmd(args []string) error {
	if len(args) <= 0 {
		return errors.New("No mac address specified to wake command")
	}

	// bcastInterface can be "eth0", "eth1", etc.. An empty string implies
	// that we use the default interface when sending the UDP packet (nil).
	bcastInterface := ""
	macAddr := args[0]

	// First we need to see if this macAddr is actually an alias, if it is:
	// we set the eth interface based on the stored item, and set the macAddr
	// based on the alias of the entry.

	// Always use the interface specified in the command line, if it exists.
	if cliFlags.BroadcastInterface != "" {
		bcastInterface = cliFlags.BroadcastInterface
	}

	// Populate the local address in the event that the broadcast interface has
	// been set.
	var localAddr *net.UDPAddr
	var err error
	if bcastInterface != "" {
		localAddr, err = ipFromInterface(bcastInterface)
		if err != nil {
			return err
		}
	}

	// The address to broadcast to is usually the default `255.255.255.255` but
	// can be overloaded by specifying an override in the CLI arguments.
	bcastAddr := fmt.Sprintf("%s:%s", cliFlags.BroadcastIP, cliFlags.UDPPort)
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}

	// Build the magic packet.
	mp, err := New(macAddr)
	if err != nil {
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	fmt.Printf("... Broadcasting to: %s\n", bcastAddr)
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Magic packet sent successfully to %s\n", macAddr)
	return nil
}
