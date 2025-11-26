package externaldns

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/rs/zerolog/log"
	"sigs.k8s.io/external-dns/endpoint"
)

const (
	identifierLabelWeight   = "weight"
	identifierLabelPriority = "priority"
)

// relativeName returns the name part of a DNS name in a zone.
func relativeName(name, zone string) string {
	return libdns.RelativeName(name, zone)
}

// absoluteName returns the FQDN of a DNS name in a zone.
func absoluteName(name, zone string) string {
	return libdns.AbsoluteName(name, zone)
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

func toExternalDNSEndpoint(record libdns.Record, zone string) *endpoint.Endpoint {
	rr := record.RR()

	endpoint := &endpoint.Endpoint{
		DNSName:    strings.TrimSuffix(absoluteName(rr.Name, zone), "."),
		Targets:    []string{rr.Data},
		RecordType: rr.Type,
		Labels:     map[string]string{},
		RecordTTL:  endpoint.TTL(rr.TTL.Seconds()),
	}

	return endpoint
}

func toLibdnsRecord(endpoint *endpoint.Endpoint, zone string) libdns.Record {
	record, err := libdns.RR{
		Type: endpoint.RecordType,
		Name: relativeName(endpoint.DNSName, zone),
		Data: endpoint.Targets[0],
		TTL:  time.Duration(endpoint.RecordTTL) * time.Second,
	}.Parse()
	if err != nil {
		log.Err(err).Msg("Failed to parse record")

		return record
	}

	return record
}

func endpointsToLibdnsZoneRecords(endpoints []*endpoint.Endpoint, zones []string) (map[string][]libdns.Record, error) {
	zoneRecords := map[string][]libdns.Record{}

	for _, endpoint := range endpoints {
		_, zone := splitDNSName(endpoint.DNSName, zones)
		if zone == "" {
			return nil, fmt.Errorf("no matching zone found for %s", endpoint.DNSName)
		}

		record := toLibdnsRecord(endpoint, zone)

		if _, ok := zoneRecords[zone]; !ok {
			zoneRecords[zone] = []libdns.Record{}
		}

		zoneRecords[zone] = append(zoneRecords[zone], record)
	}

	return zoneRecords, nil
}
