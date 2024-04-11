package pkg

import "net"

func IsIP(ip string) bool {
	addr := net.ParseIP(ip)
	if addr == nil {
		return false
	}
	return addr.To4() != nil || addr.To16() != nil
}
