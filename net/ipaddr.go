package net

import (
	"github.com/pkg/errors"
	"net"
)

func FindFirstNonLoopBackIPAddress() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve network interfaces")
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 {
			addresses, err := iface.Addrs()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to retrieve address from network interface %+v", iface)
			}
			for _, address := range addresses {
				switch addressType := address.(type) {
				case *net.IPNet:
					if !addressType.IP.IsLoopback() && addressType.IP.To4() != nil {
						return addressType.IP, nil
					}
				case *net.IPAddr:
					if !addressType.IP.IsLoopback() && addressType.IP.To4() != nil {
						return addressType.IP, nil
					}
				}
			}
		}
	}
	return nil, errors.New("failed to find an IP Address that is not a loopBack")
}
