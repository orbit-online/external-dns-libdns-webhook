package provider

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/rs/zerolog/log"
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

func toExternalDNSEndpoint(record libdns.Record, zone string) *externaldns.Endpoint {
	rr := record.RR()
	endpoint := externaldns.NewEndpointWithTTL(absoluteName(rr.Name, zone), rr.Type, int64(rr.TTL.Seconds()), rr.Data)

	switch rec := record.(type) {
	case libdns.MX:
		endpoint.WithProviderSpecific(identifierLabelPriority, strconv.FormatUint(uint64(rec.Preference), 10))
	case libdns.SRV:
		endpoint.WithProviderSpecific(identifierLabelWeight, strconv.FormatUint(uint64(rec.Weight), 10))
		endpoint.WithProviderSpecific(identifierLabelPriority, strconv.FormatUint(uint64(rec.Priority), 10))
	case libdns.ServiceBinding:
		endpoint.WithProviderSpecific(identifierLabelPriority, strconv.FormatUint(uint64(rec.Priority), 10))
	}

	return endpoint
}

func toLibdnsRecord(endpoint *externaldns.Endpoint, zone string) libdns.Record {
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
	var weight uint16
	if prop, ok := endpoint.GetProviderSpecificProperty(identifierLabelWeight); ok {
		w, err := strconv.ParseUint(prop, 10, 16)
		if err != nil {
			log.Err(err).Str("weigth", prop).Msg("Failed to parse weight")
		} else {
			weight = uint16(w)
		}
	}

	var prio uint16
	if prop, ok := endpoint.GetProviderSpecificProperty(identifierLabelPriority); ok {
		p, err := strconv.ParseUint(prop, 10, 16)
		if err != nil {
			log.Err(err).Str("priority", prop).Msg("Failed to parse priority")
		} else {
			prio = uint16(p)
		}
	}

	switch rec := record.(type) {
	case libdns.MX:
		rec.Preference = uint16(prio)
		return rec
	case libdns.SRV:
		rec.Priority = uint16(prio)
		rec.Weight = uint16(weight)
		return rec
	case libdns.ServiceBinding:
		rec.Priority = uint16(prio)
		return rec
	}

	return record
}

func endpointsToLibdnsZoneRecords(endpoints []*externaldns.Endpoint, zones []string) (map[string][]libdns.Record, error) {
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
