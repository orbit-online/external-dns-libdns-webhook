package externaldns

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/rs/zerolog/log"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

type WebhookProvider struct {
	provider.BaseProvider
	zones              []string
	libdnsProvider     libdnsregistry.Provider
	domainFilter       *endpoint.DomainFilter
	cachedZonesRecords map[string][]libdns.Record
}

func NewWebhookProvider(zones []string, libdnsProvider libdnsregistry.Provider) *WebhookProvider {
	return &WebhookProvider{
		zones:              zones,
		domainFilter:       endpoint.NewDomainFilter(zones),
		libdnsProvider:     libdnsProvider,
		cachedZonesRecords: map[string][]libdns.Record{},
	}
}

func (p WebhookProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	endpoints := []*endpoint.Endpoint{}

	// return all records for configured zones
	for _, zone := range p.zones {
		logger := log.With().Str("zone", zone).Logger()
		logger.Debug().Msg("Retrieving records for zone")

		records, err := p.libdnsProvider.GetRecords(ctx, zone)
		if err != nil {
			logger.Err(err).Msg("Failed to retrieve records for zone")

			return nil, fmt.Errorf("failed to retrieve records for zone: %w", err)
		}

		// as there is no real concurrent sync in progress we can cache between the calls to avoid calling api too many times
		p.cachedZonesRecords[zone] = records

		for _, record := range records {
			rr := record.RR()
			endpoint := &endpoint.Endpoint{
				DNSName:    strings.TrimSuffix(libdns.AbsoluteName(rr.Name, zone), "."),
				Targets:    []string{rr.Data},
				RecordType: rr.Type,
				Labels:     map[string]string{},
				RecordTTL:  endpoint.TTL(rr.TTL.Seconds()),
			}
			logger.Trace().
				Any("record", record).
				Any("endpoint", endpoint).
				Msg("Record converted to endpoint")

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

func (p WebhookProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	log.Debug().
		Any("changes_create", changes.Create).
		Msg("Convert creation change endpoints to records")

	creates, err := endpointsToLibdnsZoneRecords(changes.Create, p.zones)
	if err != nil {
		return err
	}

	log.Debug().
		Any("changes_delete", changes.Delete).
		Msg("Convert deletion change endpoints to records")

	deletes, err := endpointsToLibdnsZoneRecords(changes.Delete, p.zones)
	if err != nil {
		return err
	}

	log.Debug().
		Any("changes_update_new", changes.UpdateNew).
		Any("changes_update_old", changes.UpdateOld).
		Msg("Convert updates change endpoints to records")

	updates, err := endpointsToLibdnsZoneRecords(changes.UpdateNew, p.zones)
	if err != nil {
		return err
	}

	if len(creates) > 0 {
		for zone, records := range creates {
			log.Info().
				Any("records", records).
				Any("zone", zone).
				Msg("Creating records")

			_, err := p.libdnsProvider.AppendRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to create records: %w", err)
			}
		}
	}

	if len(deletes) > 0 {
		for zone, records := range deletes {
			log.Info().
				Any("records", records).
				Any("zone", zone).
				Msg("Deleting records")

			_, err := p.libdnsProvider.DeleteRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to delete records: %w", err)
			}
		}
	}

	if len(updates) > 0 {
		for zone, records := range updates {
			log.Info().
				Any("records", records).
				Any("zone", zone).
				Msg("Updating records")

			_, err := p.libdnsProvider.SetRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to update records: %w", err)
			}
		}
	}

	return nil
}

func endpointsToLibdnsZoneRecords(endpoints []*endpoint.Endpoint, zones []string) (map[string][]libdns.Record, error) {
	zoneRecords := map[string][]libdns.Record{}

	for _, endpoint := range endpoints {
		_, zone := splitDNSName(endpoint.DNSName, zones)
		if zone == "" {
			return nil, fmt.Errorf("no matching zone found for %s", endpoint.DNSName)
		}

		record := libdns.RR{
			Type: endpoint.RecordType,
			Name: libdns.RelativeName(endpoint.DNSName, zone),
			Data: endpoint.Targets[0],
			TTL:  time.Duration(endpoint.RecordTTL) * time.Second,
		}

		if _, ok := zoneRecords[zone]; !ok {
			zoneRecords[zone] = []libdns.Record{}
		}

		zoneRecords[zone] = append(zoneRecords[zone], record)
	}

	return zoneRecords, nil
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
