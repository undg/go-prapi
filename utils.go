package main

import (
	"log"
	"net"
	"net/http"
	"strings"
)

func actionsToStrings(actions []Action) []string {
	strs := make([]string, len(actions))
	for i, action := range actions {
		strs[i] = string(action)
	}
	return strs
}

func IsLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}

	return false
}

func upgraderCheckOrigin() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error splitting host and port: %v", err)
			return false
		}

		ip := net.ParseIP(host)
		if ip == nil {
			log.Printf("Invalid IP: %s", host)
			return false
		}

		return IsLocalIP(ip) || strings.HasPrefix(r.Host, "localhost")
	}
}

