package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	hierr "github.com/reconquest/hierr-go"
)

func discovery() error {
	discoveryData := make(map[string][]map[string]string)

	out, err := exec.Command(sudo, smartctl, "--scan-open").CombinedOutput()
	if err != nil {
		return hierr.Errorf(
			out,
			"can't run %s %s --scan-open",
			sudo,
			smartctl,
		)
	}

	disks, err := getDisks(string(out))
	if err != nil {
		return hierr.Errorf(err, "can't discovery disks")
	}

	for _, diskData := range disks {
		discoveryData["data"] = append(discoveryData["data"], diskData)
	}

	out, err = json.Marshal(discoveryData)
	if err != nil {
		return hierr.Errorf(err, "can't create json data")
	}

	fmt.Printf("%s\n", out)

	return nil
}

func getDisks(rawOut string) ([]map[string]string, error) {
	var disks []map[string]string
	regexpDiskName := regexp.MustCompile(`(\/(.+?))\s(-d [A-Za-z0-9,\+]+)`)

	scanner := bufio.NewScanner(strings.NewReader(rawOut))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Permission denied") {
			return nil, fmt.Errorf("can't get disks. Permission denied")
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		diskName := regexpDiskName.FindString(line)

		diskData, err := getDiskData(diskName)
		if err != nil {
			return nil, hierr.Errorf(
				err,
				"can't get data for %s disk",
				diskName,
			)
		}
		disks = append(disks, diskData)
	}

	return disks, nil
}

func getDiskData(diskName string) (map[string]string, error) {
	diskData := make(map[string]string)

	diskData["{#DISK_NAME}"] = diskName
	diskData["{#SMART_ENABLED}"] = "Disabled"
	diskData["{#INTERFACE_TYPE}"] = SAS

	args := []string{smartctl, "-i"}
	args = append(args, strings.Split(diskName, " ")...)

	rawDiskInfo, err := exec.Command(sudo, args...).CombinedOutput()
	if err != nil {
		return diskData, hierr.Errorf(
			rawDiskInfo,
			"can't run %s %s -i %s",
			sudo,
			smartctl,
			diskName,
		)
	}

	diskInfo := string(rawDiskInfo)

	r := regexp.MustCompile(`SMART.+?: +(.+)`)
	smartStats := r.FindAllString(diskInfo, -1)
	for _, smartStat := range smartStats {
		if strings.Contains(smartStat, "Enabled") {
			diskData["{#SMART_ENABLED}"] = "Enabled"
		}
	}

	if strings.Contains(diskInfo, SATA) {
		diskData["{#INTERFACE_TYPE}"] = SATA
	}

	return diskData, nil
}
