package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	tls "github.com/hellais/utls-light/tls"
)

func getRequest(conn net.Conn, requestHostname string, alpn string) (*http.Response, error) {
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Host: requestHostname},
		Header: make(http.Header),
		Host:   requestHostname,
	}
	err := req.Write(conn)
	if err != nil {
		return nil, err
	}
	return http.ReadResponse(bufio.NewReader(conn), req)
}

func logStatus(status string, serverName string, addr string) {
	f, err := os.OpenFile("domain-status.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	logline := fmt.Sprintf("%s,%s,%s\n", status, serverName, addr)
	fmt.Print(logline)
	if _, err = f.WriteString(logline); err != nil {
		panic(err)
	}
}

func testURL(serverName string) error {
	ips, err := net.LookupIP(serverName)
	if len(ips) == 0 || err != nil {
		logStatus("FAIL-DNS", serverName, "")
		return errors.New("failed to lookup IP")
	}
	addr := fmt.Sprintf("%s:443", ips[0].String())
	config := tls.Config{ServerName: serverName}
	dialConn, err := net.DialTimeout("tcp", addr, time.Duration(2)*time.Second)
	if err != nil {
		logStatus("FAIL-CONNECT", serverName, addr)
		return err
	}
	tlsConn := tls.Client(dialConn, &config)
	defer tlsConn.Close()
	_, err = getRequest(tlsConn, serverName, tlsConn.ConnectionState().NegotiatedProtocol)
	if err != nil {
		logStatus("FAIL-GET", serverName, addr)
		return err
	}
	logStatus("OK", serverName, addr)
	return nil
}

func main() {
	file, err := os.Open("tests/citizenlab-domains.txt")
	if err != nil {
		log.Fatal("Cannot open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := scanner.Text()
		err := testURL(domain)
		if err != nil {
			log.Printf("Failed to check %s %v", domain, err)
		}
	}

}
