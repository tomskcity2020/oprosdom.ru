package shared_validate

import (
	"errors"
	"net"
)

func IpValidate(p string) (net.IP, error) {
	ip := net.ParseIP(p)
	if ip == nil {
		return nil, errors.New("ip_not_valid")
	}

	return ip, nil
}
