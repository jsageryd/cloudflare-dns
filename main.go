package main

import (
	"fmt"
	"log/syslog"
	"os"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

func main() {
	domains := os.Args[1:]
	if len(domains) == 0 {
		fmt.Println("no domains listed")
		os.Exit(1)
	}
	apiKey := os.Getenv("CF_API_KEY")
	if apiKey == "" {
		fmt.Println("CF_API_KEY not set")
		os.Exit(1)
	}
	apiEmail := os.Getenv("CF_API_EMAIL")
	if apiEmail == "" {
		fmt.Println("CF_API_EMAIL not set")
		os.Exit(1)
	}

	extIPF := extIPFuture(10 * time.Second)

	cf, err := NewCF(apiKey, apiEmail)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fetchDNSARecordsF := cf.fetchDNSARecordsFuture(domains...)

	sl, err := syslog.New(0, "cloudflare-dns")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ip, err := extIPF()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dnsARecords, err := fetchDNSARecordsF()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var forUpdate []cloudflare.DNSRecord
	for _, r := range dnsARecords {
		if ip != r.Content {
			forUpdate = append(forUpdate, r)
			continue
		}
		sl.Info(fmt.Sprintf("[%s] no change\n", r.ZoneName))
	}

	for _, r := range forUpdate {
		var oldIP string
		oldIP, r.Content = r.Content, ip
		err := cf.updateDNSRecord(r)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sl.Info(fmt.Sprintf("[%s] %s -> %s\n", r.ZoneName, oldIP, r.Content))
	}
}
