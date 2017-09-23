package controllers

import (
	"math/rand"
	"strings"

	"clickyab.com/cluster-tools/modules/k8s/config"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/assert"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

var availableIPs []string

// lint bypass only
func init() {
	dnsRefresher([]string{}, map[string][]string{})
}

// zones is smth like {clickyab.com: {"www", "v", "a"}}
func dnsRefresher(ips []string, zones map[string][]string) {
	availableIPs = ips
	api, err := cloudflare.New(kcfg.CFKey, kcfg.CFMail)
	assert.Nil(err)

	// loop over domains
	for domain := range zones {
		zoneID, err := api.ZoneIDByName(domain)
		if err != nil {
			continue
		}

		// loop over subdomains of a domain
		for _, subDomain := range zones[domain] {
			fullDomain := subDomain + "." + domain
			ipToRecord := getRecordsDetail(api, zoneID, fullDomain)
			clearDownIPs(api, ipToRecord, availableIPs)

			newIPToRecord := getRecordsDetail(api, zoneID, fullDomain)
			if len(newIPToRecord) < kcfg.LeastIPNum {
				refreshIPs(api, newIPToRecord, fullDomain)
			}
		}
	}
}

func getRecordsDetail(api *cloudflare.API, zoneID, domain string) map[string]cloudflare.DNSRecord {
	records, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: domain})
	assert.Nil(err)

	var ipToRecord = make(map[string]cloudflare.DNSRecord)
	for _, record := range records {
		if record.Type != "A" {
			continue
		}
		ipToRecord[record.Content] = record
	}

	return ipToRecord
}

func clearDownIPs(api *cloudflare.API, ipMap map[string]cloudflare.DNSRecord, ips []string) {
	for ip := range ipMap {
		if !array.StringInArray(ip, ips...) {
			err := api.DeleteDNSRecord(ipMap[ip].ZoneID, ipMap[ip].ID)
			assert.Nil(err)
		}
	}
}

func refreshIPs(api *cloudflare.API, ipToRecordMap map[string]cloudflare.DNSRecord, domain string) {
	var expected int = kcfg.LeastIPNum
	if len(availableIPs) < kcfg.LeastIPNum {
		expected = len(availableIPs)
	}

	targetNum := expected - len(availableIPs)

	for targetNum != 0 {
		ip := availableIPs[rand.Intn(len(availableIPs))]
		if _, ok := ipToRecordMap[ip]; ok {
			continue
		}

		zoneNameArray := strings.Split(domain, ".")[1:]
		zoneName := strings.Join(zoneNameArray, ".")

		zoneID, err := api.ZoneIDByName(zoneName)
		assert.Nil(err)

		if _, err = api.CreateDNSRecord(zoneID, cloudflare.DNSRecord{
			Name:    domain,
			Type:    "A",
			ZoneID:  zoneID,
			Content: ip,
			Proxied: true,
		}); err != nil {
			continue
		}
		targetNum--
	}
}
