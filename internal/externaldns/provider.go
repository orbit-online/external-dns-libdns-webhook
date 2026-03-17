package externaldns

import (
	"context"
	"errors"
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
	var errs []error
	endpoints := []*endpoint.Endpoint{}

	// return all records for configured zones
	for _, zone := range p.domainFilter.Filters {
		slog.Debug("Retrieving records for zone")

		records, err := p.libdnsProvider.GetRecords(ctx, zone)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to retrieve records for zone %s: %w", zone, err))
			continue
		}

		groupedRecords := map[string](map[string][]libdns.Record){}
		for _, record := range records {
			rr := record.RR()
			if _, ok := groupedRecords[rr.Type]; !ok {
				groupedRecords[rr.Type] = map[string][]libdns.Record{}
			}
			groupedRecords[rr.Type][rr.Name] = append(groupedRecords[rr.Type][rr.Name], record)
		}
		for recordType, typeGroup := range groupedRecords {
			for recordName, nameGroup := range typeGroup {
				var data []string
				for _, record := range nameGroup {
					data = append(data, record.RR().Data)
				}
				ep := &endpoint.Endpoint{
					DNSName:    strings.TrimSuffix(libdns.AbsoluteName(recordName, zone), "."),
					Targets:    data,
					RecordType: recordType,
					Labels:     map[string]string{},
					RecordTTL:  endpoint.TTL(nameGroup[0].RR().TTL.Seconds()),
				}
				slog.Debug("Converted records to endpoint", "records", nameGroup, "endpoint", ep)

				endpoints = append(endpoints, ep)
			}
		}
	}

	return endpoints, errors.Join(errs...)
}

func (p WebhookProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	var errs []error

	endpointsToLibdnsZoneRecords := func(endpoints []*endpoint.Endpoint) map[string][]libdns.Record {
		zoneRecords := map[string][]libdns.Record{}
		for _, ep := range endpoints {
			_, zone := splitDNSName(ep.DNSName, p.domainFilter.Filters)
			if zone == "" {
				slog.Debug("no matching zone found for endpoint", "ep", ep)
				continue
			}
			for _, target := range ep.Targets {
				record, err := libdns.RR{
					Type: ep.RecordType,
					Name: libdns.RelativeName(ep.DNSName, zone),
					Data: strings.Trim(target, "\""),
					TTL:  time.Duration(ep.RecordTTL) * time.Second,
				}.RR().Parse()
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse endpoint target %s, endpoint: %+v, err: %w", target, ep, err))
					continue
				}
				slog.Debug("Converted endpoint to record", "endpoint", ep, "record", record)

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
			created, err := p.libdnsProvider.AppendRecords(ctx, zone, records)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to create records in zone %s: %w", zone, err))
			} else {
				if len(created) != len(records) {
					errs = append(errs, fmt.Errorf("number of created records (%d) did not match number of records to create (%d)", len(created), len(records)))
				} else {
					slog.Debug("records created", "actual", len(created), "wanted", len(records))
				}
			}
		}
	}

	if len(deletes) > 0 {
		for zone, records := range deletes {
			slog.Info("Deleting records", "zone", zone, "records", records)
			deleted, err := p.libdnsProvider.DeleteRecords(ctx, zone, records)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to delete records in zone %s: %w", zone, err))
			} else {
				if len(deleted) != len(records) {
					errs = append(errs, fmt.Errorf("number of deleted records (%d) did not match number of records to delete (%d)", len(deleted), len(records)))
				} else {
					slog.Debug("records deleted", "actual", len(deleted), "wanted", len(records))
				}
			}
		}
	}

	if len(updates) > 0 {
		for zone, records := range updates {
			slog.Info("Updating records", "zone", zone, "records", records)
			updated, err := p.libdnsProvider.SetRecords(ctx, zone, records)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to update records in zone %s: %w", zone, err))
			} else {
				if len(updated) != len(records) {
					errs = append(errs, fmt.Errorf("number of updated records (%d) did not match number of records to update (%d)", len(updated), len(records)))
				} else {
					slog.Debug("records updated", "actual", len(updated), "wanted", len(records))
				}
			}
		}
	}

	return errors.Join(errs...)
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
