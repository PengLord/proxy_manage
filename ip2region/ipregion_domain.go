package ip2region

import (
	"strings"
	"os/exec"
	"strconv"
	"log"
)

var region *Ip2Region

func InitRegion(path string) {
	var err error
	region, err = New(path)
	if err != nil {
		log.Fatal(err)
	}
}

func IsCn(domain string) bool {
	addressMap := nsLookup(domain)
	for k, v := range addressMap {
		if k == "address1" {
			ipInfo, _ := region.MemorySearch(v)
			if ipInfo.Country == "中国" {
				return true
			}
			continue
		}

	}
	return false

}

func nsLookup(address string) map[string]string {
	canonicalData := make(map[string]string)
	canonicalCount, nameCount, addressCount := 1, 1, 1
	out, _ := exec.Command("nslookup", address).Output()
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[2:] {
		if strings.Contains(line, "Name:") {
			canonicalData["name"+strconv.Itoa(nameCount)] = strings.Fields(line)[1]
			nameCount++
		} else if strings.Contains(line, "Address:") {
			canonicalData["address"+strconv.Itoa(addressCount)] = strings.Fields(line)[1]
			addressCount++
		} else if strings.Contains(line, "canonical") {
			cname := strings.Fields(line)
			canonicalData["cname"+strconv.Itoa(canonicalCount)] = cname[len(cname)-1]
			canonicalCount++
		}
	}
	return canonicalData
}
