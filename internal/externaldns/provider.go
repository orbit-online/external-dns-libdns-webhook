package externaldns

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

type WebhookProvider struct {
	provider.BaseProvider
	domainFilter   *endpoint.DomainFilter
	libdnsProvider libdnsregistry.Provider
}

func NewWebhookProvider(zones []string, libdnsProvider libdnsregistry.Provider) *WebhookProvider {
	return &WebhookProvider{
		domainFilter:   endpoint.NewDomainFilter(zones),
		libdnsProvider: libdnsProvider,
	}
}

func (p WebhookProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	errs := 0
	endpoints := []*endpoint.Endpoint{}

	// return all records for configured zones
	for _, zone := range p.domainFilter.Filters {
		slog.Debug("Retrieving records for zone")

		records, err := p.libdnsProvider.GetRecords(ctx, zone)
		if err != nil {
			errs++
			slog.Error("failed to retrieve records for zone", "zone", zone, "err", err)
			continue
		}

		for _, record := range records {
			rr := record.RR()
			endpoint := &endpoint.Endpoint{
				DNSName:    strings.TrimSuffix(libdns.AbsoluteName(rr.Name, zone), "."),
				Targets:    []string{rr.Data},
				RecordType: rr.Type,
				Labels:     map[string]string{},
				RecordTTL:  endpoint.TTL(rr.TTL.Seconds()),
			}
			slog.Debug("Converted record to endpoint", "record", record, "endpoint", endpoint)

			endpoints = append(endpoints, endpoint)
		}
	}

	if errs > 0 {
		return endpoints, fmt.Errorf("encountered %d errors while retrieving records", errs)
	} else {
		return endpoints, nil
	}
}

func (p WebhookProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	errs := 0

	endpointsToLibdnsZoneRecords := func(endpoints []*endpoint.Endpoint) map[string][]libdns.Record {
		zoneRecords := map[string][]libdns.Record{}

		for _, endpoint := range endpoints {
			_, zone := splitDNSName(endpoint.DNSName, p.domainFilter.Filters)
			if zone == "" {
				errs++
				slog.Error("no matching zone found for endpoint", "endpoint", endpoint)
			} else {
				record := libdns.RR{
					Type: endpoint.RecordType,
					Name: libdns.RelativeName(endpoint.DNSName, zone),
					Data: endpoint.Targets[0],
					TTL:  time.Duration(endpoint.RecordTTL) * time.Second,
				}
				slog.Debug("Converted endpoint to record", "endpoint", endpoint, "record", record)

				if _, ok := zoneRecords[zone]; !ok {
					zoneRecords[zone] = []libdns.Record{}
				}
				zoneRecords[zone] = append(zoneRecords[zone], record)
			}
		}

		return zoneRecords
	}

	slog.Debug("Converting changes.Create to records", "changes.Create", changes.Create)
	creates := endpointsToLibdnsZoneRecords(changes.Create)
	slog.Debug("Converting changes.Delete to records", "changes.Delete", changes.Delete)
	deletes := endpointsToLibdnsZoneRecords(changes.Delete)
	slog.Debug("Converting changes.UpdateOld/New to records", "changes.UpdateOld", changes.UpdateOld, "changes.UpdateNew", changes.UpdateNew)
	updates := endpointsToLibdnsZoneRecords(changes.UpdateNew)

	if len(creates) > 0 {
		for zone, records := range creates {
			slog.Info("Creating records", "zone", zone, "records", records)
			_, err := p.libdnsProvider.AppendRecords(ctx, zone, records)
			if err != nil {
				errs++
				slog.Error("failed to create records", "err", err)
			}
		}
	}

	if len(deletes) > 0 {
		for zone, records := range deletes {
			slog.Info("Deleting records", "zone", zone, "records", records)
			_, err := p.libdnsProvider.DeleteRecords(ctx, zone, records)
			if err != nil {
				errs++
				slog.Error("failed to delete records", "err", err)
			}
		}
	}

	if len(updates) > 0 {
		for zone, records := range updates {
			slog.Info("Updating records", "zone", zone, "records", records)
			_, err := p.libdnsProvider.SetRecords(ctx, zone, records)
			if err != nil {
				errs++
				slog.Error("failed to update records", "err", err)
			}
		}
	}

	if errs > 0 {
		return fmt.Errorf("encountered %d errors while applying changes", errs)
	} else {
		return nil
	}
}

// splitDNSName splits a DNS name into a name and a zone.
func splitDNSName(dnsName string, zones []string) (string, string) {
	name := strings.TrimSuffix(dnsName, ".")

	domain := ""
	resourceRecord := ""
	// sort zones by dot count; make sure subdomains sort earlier
	sort.Slice(zones, func(i, j int) bool {
		return strings.Count(zones[i], ".") > strings.Count(zones[j], ".")
	})

	for _, filter := range zones {
		if strings.HasSuffix(name, "."+filter) {
			domain = filter
			resourceRecord = name[0 : len(name)-len(filter)-1]

			break
		} else if name == filter {
			domain = filter
			resourceRecord = ""
		}
	}

	return resourceRecord, domain
}
