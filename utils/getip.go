package utils

import (
	"fmt"
	"httphere/conf"
	"net"
)

func GetLanIp() (string, string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	ipv4 := ""
	ipv6 := ""

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipv4 = ipnet.IP.String()
			}
			if ipnet.IP.To16() != nil {
				ipv6 = ipnet.IP.String()
			}
		}
	}

	return ipv4, ipv6
}

func GetUrls() (string, string, string) {
	ipv4, ipv6 := GetLanIp()
	ipStr := ipv6
	if ipv4 != "" {
		ipStr = ipv4
	}

	viewUrl := ""
	uploadUrl := ""
	qrUrl := ""
	if ipStr != "" {
		viewUrl = fmt.Sprintf("%v%v:%v", "http://", ipStr, conf.GetPort())

		uploadUrl = fmt.Sprintf("%v%v:%v/httphere_upload", "http://", ipStr, conf.GetPort())

		qrUrl = fmt.Sprintf("%v%v:%v/qr", "http://", ipStr, conf.GetPort())
	}

	return viewUrl, uploadUrl, qrUrl
}
