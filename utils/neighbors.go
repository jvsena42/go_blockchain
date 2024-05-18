package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

func IsFoundNode(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("IsFoundNode error: %s %v \n", target, err)
		return false
	}

	return true
}

var IP_PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbors(myHostIP string, myPort uint16, startIP uint8, endIP uint8, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", myHostIP, myPort)

	m := IP_PATTERN.FindStringSubmatch(myHostIP)
	if m == nil {
		return nil
	}
	ipPrefix := m[1]
	hostIdent, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	for guessPort := startPort; guessPort <= endPort; guessPort++ {
		for variableHostIndent := startIP; variableHostIndent <= endIP; variableHostIndent++ {
			guessIP := fmt.Sprintf("%s%d", ipPrefix, hostIdent+int(variableHostIndent))
			guessTarget := fmt.Sprintf("%s:%d", guessIP, guessPort)
			if guessTarget != address && IsFoundNode(guessIP, guessPort) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	return neighbors
}

func GetHost() string {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	address := conn.LocalAddr().(*net.UDPAddr)
	ipStr := fmt.Sprintf("%v", address.IP)

	return ipStr
}
