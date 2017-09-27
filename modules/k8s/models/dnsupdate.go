package models

import (
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

var (
	cfKey      = config.RegisterString("cmon.modules.cloudflare.key", "EMPTY", "cloudflare key")
	cfMail     = config.RegisterString("cmon.modules.cloudflare.mail", "EMPTY", "cloudflare mail")
	leastIPNum = config.RegisterInt("cmon.modules.cloudflare.min_ip_per_dns", 3, "minimum ip per dns record")
)

// RefreshDNS is the function to handle records in cloudflare
// zones is some thing like {clickyab.com: {"www", "v", "a", "@"}}
func RefreshDNS(ips []string, zones map[string][]string) {
	api, err := cloudflare.New(cfKey.String(), cfMail.String())
	assert.Nil(err)

	// loop over domains
	for domain := range zones {
		zoneID, err := api.ZoneIDByName(domain)
		if err != nil {
			logrus.WithError(err).WithField("domain", domain).Debugf("can not get zone id")
			continue
		}

		// loop over sub domains of a domain
		for _, subDomain := range zones[domain] {
			fullDomain := domain
			if subDomain != "@" {
				fullDomain = subDomain + "." + fullDomain
			}
			ipToRecord := getRecordsDetail(api, zoneID, fullDomain)
			cnt := clearDownIPs(api, ipToRecord, ips)
			if cnt < leastIPNum.Int() {
				refreshIPs(api, ipToRecord, zoneID, fullDomain, ips...)
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

func clearDownIPs(api *cloudflare.API, ipMap map[string]cloudflare.DNSRecord, ips []string) int {
	i := len(ipMap)
	for ip := range ipMap {
		if !array.StringInArray(ip, ips...) {
			i--
			err := api.DeleteDNSRecord(ipMap[ip].ZoneID, ipMap[ip].ID)
			assert.Nil(err)
			delete(ipMap, ip)
		}
	}

	return i
}

func refreshIPs(api *cloudflare.API, ipToRecordMap map[string]cloudflare.DNSRecord, zoneID, domain string, ips ...string) {
	logrus.WithField("domain", domain).Debugf("refreshing ips")
	availableIPs := make(map[string]bool)
	for i := range ips {
		availableIPs[ips[i]] = true
	}
	var expected int = leastIPNum.Int()
	if len(availableIPs) < expected {
		logrus.Warn("les than requested ips")
		expected = len(availableIPs)
	}
	expected -= len(ipToRecordMap)
	for ip := range availableIPs {
		if expected <= 0 {
			break
		}
		if _, err := api.CreateDNSRecord(zoneID, cloudflare.DNSRecord{
			Name:    domain,
			Type:    "A",
			ZoneID:  zoneID,
			Content: ip,
			Proxied: true,
		}); err != nil {
			logrus.WithError(err).WithField("domain", domain).Debug("can not create record")
			continue
		}
		expected--
	}
}
