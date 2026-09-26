package main

import (
	"context"
	"errors"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

// aRecord is a DNS A record along with its zone. The Cloudflare API no longer
// includes zone_id and zone_name in DNS record responses.
type aRecord struct {
	ZoneID   string
	ZoneName string
	cloudflare.DNSRecord
}

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

func (c *cfClient) fetchDNSARecordsFuture(ctx context.Context, domains ...string) func() ([]aRecord, error) {
	var dnsRecords []aRecord
	var err error

	done := make(chan struct{})
	go func() {
		dnsRecords, err = c.fetchDNSARecords(ctx, domains...)
		close(done)
	}()

	return func() ([]aRecord, error) {
		<-done
		return dnsRecords, err
	}
}

func (c *cfClient) fetchDNSARecords(ctx context.Context, domains ...string) ([]aRecord, error) {
	zones, err := c.cf.ListZones(ctx, domains...)
	if err != nil {
		return nil, err
	}

	if len(zones) == 0 {
		return nil, errors.New("no matching domains")
	}

	var dnsRecords []aRecord
	for _, z := range zones {
		drs, _, err := c.cf.ListDNSRecords(
			ctx,
			cloudflare.ZoneIdentifier(z.ID),
			cloudflare.ListDNSRecordsParams{
				Type: "A",
			},
		)
		if err != nil {
			return nil, err
		}
		for _, dr := range drs {
			dnsRecords = append(dnsRecords, aRecord{
				ZoneID:    z.ID,
				ZoneName:  z.Name,
				DNSRecord: dr,
			})
		}
	}

	return dnsRecords, nil
}

func (c *cfClient) updateDNSRecord(ctx context.Context, r aRecord) error {
	_, err := c.cf.UpdateDNSRecord(
		ctx,
		cloudflare.ZoneIdentifier(r.ZoneID),
		cloudflare.UpdateDNSRecordParams{
			ID:      r.ID,
			Type:    r.Type,
			Name:    r.Name,
			Content: r.Content,
			Tags:    r.Tags,
		},
	)
	return err
}
