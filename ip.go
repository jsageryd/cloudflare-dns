package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func extIPFuture(timeout time.Duration) func() (string, error) {
	var ip string
	var err error

	done := make(chan struct{})
	go func() {
		ip, err = extIP(timeout)
		close(done)
	}()

	return func() (string, error) {
		<-done
		return ip, err
	}
}

func extIP(timeout time.Duration) (string, error) {
	hc := &http.Client{Timeout: timeout}
	res, err := hc.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	lr := io.LimitReader(res.Body, 256)
	ipBytes, err := io.ReadAll(lr)
	if err != nil {
		return "", err
	}
	ip := string(bytes.TrimRight(ipBytes, "\n"))
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP: %s", ip)
	}
	return ip, nil
}
