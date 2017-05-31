package ipcalc

import (
	"errors"
	"math/big"
	"net"
	"strings"
)

var invalidRangeError = errors.New("Invalid IP Range")

func GetIPsFromCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if ip[3] != 0 && ip[3] != 255 {
			ips = append(ips, ip.String())
		}
	}

	return ips, nil
}

func GetIPsFromRange(ipRange string) ([]string, error) {
	firstAndLastIPs := strings.Split(ipRange, "-")
	if len(firstAndLastIPs) == 1 {
		if net.ParseIP(ipRange) == nil {
			return nil, invalidRangeError
		}
		return []string{ipRange}, nil
	}
	return listIPsInRange(firstAndLastIPs[0], firstAndLastIPs[1])
}

func listIPsInRange(firstIPString string, lastIPString string) ([]string, error) {
	firstIP := net.ParseIP(firstIPString)
	lastIP := net.ParseIP(lastIPString)

	if lastIP == nil || firstIP == nil || isAReversedRange(firstIP, lastIP) {
		return nil, invalidRangeError
	}

	return listIPsInSafeRange(firstIP, lastIP), nil
}

func listIPsInSafeRange(firstIP net.IP, lastIP net.IP) []string {
	var ips []string
	for ip := firstIP; !ip.Equal(lastIP); inc(ip) {
		ips = append(ips, ip.String())
	}
	return append(ips, lastIP.String())
}

func isAReversedRange(firstIP net.IP, lastIP net.IP) bool {
	firstIPAsInt := new(big.Int).SetBytes(firstIP)
	lastIPAsInt := new(big.Int).SetBytes(lastIP)
	return firstIPAsInt.Cmp(lastIPAsInt) == 1
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
