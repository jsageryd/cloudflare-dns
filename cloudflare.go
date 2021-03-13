package main

import (
	"context"
	"errors"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type cfClient struct {
	cf *cloudflare.API
}

func NewCF(apiKey, apiEmail string) (*cfClient, error) {
	cf, err := cloudflare.New(apiKey, apiEmail)
	if err != nil {
		return nil, err
	}
	return &cfClient{cf: cf}, nil
}

func (c *cfClient) fetchDNSARecordsFuture(ctx context.Context, domains ...string) func() ([]cloudflare.DNSRecord, error) {
	var dnsRecords []cloudflare.DNSRecord
	var err error

	done := make(chan struct{})
	go func() {
		dnsRecords, err = c.fetchDNSARecords(ctx, domains...)
		close(done)
	}()

	return func() ([]cloudflare.DNSRecord, error) {
		<-done
		return dnsRecords, err
	}
}

func (c *cfClient) fetchDNSARecords(ctx context.Context, domains ...string) ([]cloudflare.DNSRecord, error) {
	zones, err := c.cf.ListZones(ctx, domains...)
	if err != nil {
		return nil, err
	}

	if len(zones) == 0 {
		return nil, errors.New("no matching domains")
	}

	var dnsRecords []cloudflare.DNSRecord
	for _, z := range zones {
		drs, err := c.cf.DNSRecords(
			ctx,
			z.ID,
			cloudflare.DNSRecord{
				Type: "A",
			},
		)
		if err != nil {
			return nil, err
		}
		dnsRecords = append(dnsRecords, drs...)
	}

	return dnsRecords, nil
}

func (c *cfClient) updateDNSRecord(ctx context.Context, r cloudflare.DNSRecord) error {
	return c.cf.UpdateDNSRecord(ctx, r.ZoneID, r.ID, r)
}
