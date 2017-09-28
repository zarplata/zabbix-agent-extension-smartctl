package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	hierr "github.com/reconquest/hierr-go"
)

func getDiskStats(diskInfo map[string]string) (map[string]string, error) {
	args := []string{smartctl, "-a"}
	args = append(args, strings.Split(diskInfo["name"], " ")...)
	out, err := exec.Command(sudo, args...).CombinedOutput()
	if err != nil {
		return nil, hierr.Errorf(
			out, "can't get stats for %s disk",
			diskInfo["name"],
		)
	}

	if diskInfo["interface"] == SAS {
		return getSASStats(diskInfo, string(out))
	}

	return getSATAStats(diskInfo, string(out))
}

func getSASStats(
	diskInfo map[string]string,
	rawStats string,
) (map[string]string, error) {
	diskStats := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(rawStats))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Vendor:") {
			metric := strings.Split(line, ":")[1]
			diskStats["diskIDSAS"] = strings.TrimSpace(metric)
		}

		if strings.Contains(line, "Product:") {
			metric := strings.Split(line, ":")[1]
			diskStats["diskIDSAS"] = fmt.Sprintf(
				"%s_%s",
				diskStats["diskIDSAS"],
				strings.TrimSpace(metric),
			)
		}

		if strings.Contains(line, "Serial number:") {
			metric := strings.Split(line, ":")[1]
			diskStats["diskIDSAS"] = fmt.Sprintf(
				"%s-%s",
				diskStats["diskIDSAS"],
				strings.TrimSpace(metric),
			)
		}
		if strings.Contains(line, "Health Status") {
			metric := strings.Split(line, ":")[1]
			diskStats["healthStatus"] = strings.TrimSpace(metric)
		}

		if strings.Contains(line, "Current Drive Temperature") {
			metric := strings.Replace(strings.Split(line, ":")[1], "C", "", -1)
			diskStats["currentDriveTemp"] = strings.TrimSpace(metric)
		}

		if strings.Contains(line, "Elements in grown defect list") {
			metric := strings.Split(line, ":")[1]
			diskStats["elementsInGrownDefectList"] = strings.TrimSpace(metric)
		}

		if strings.Contains(line, "Non-medium error count") {
			metric := strings.Split(line, ":")[1]
			diskStats["non-mediumErrorCount"] = strings.TrimSpace(metric)
		}

		if strings.Contains(line, "read:") {
			metric := strings.Fields(line)[7]
			diskStats["uncorrectedErrorsRead"] = strings.TrimSpace(metric)
		}

		if strings.Contains(line, "write:") {
			metric := strings.Fields(line)[7]
			diskStats["uncorrectedErrorsWrite"] = strings.TrimSpace(metric)
		}

		if strings.Contains(line, "verify:") {
			metric := strings.Fields(line)[7]
			diskStats["uncorrectedErrorsVerify"] = strings.TrimSpace(metric)
		}
	}

	return diskStats, nil
}

func getSATAStats(
	diskInfo map[string]string,
	rawStats string,
) (map[string]string, error) {
	diskStats := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(rawStats))
	for scanner.Scan() {
		line := scanner.Text()
		attribute := strings.Fields(line)

		if strings.Contains(line, "Device Model:") {
			metric := strings.Split(line, ":")[1]
			diskStats["diskIDSATA"] = strings.Replace(
				strings.TrimSpace(metric),
				" ",
				"_",
				-1,
			)
		}

		if strings.Contains(line, "Serial Number:") {
			metric := strings.Split(line, ":")[1]
			diskStats["diskIDSATA"] = fmt.Sprintf(
				"%s-%s",
				diskStats["diskIDSATA"],
				strings.TrimSpace(metric),
			)
		}

		if len(attribute) >= 9 {
			if attribute[0] == "5" {
				diskStats["reallocatedSectorCount"] = attribute[9]
			}

			if attribute[0] == "184" {
				diskStats["endToEndError"] = attribute[9]
			}

			if attribute[0] == "187" {
				diskStats["reportedUncorrectableErrors"] = attribute[9]
			}

			if attribute[0] == "194" {
				diskStats["temperature"] = attribute[9]
			}

			if attribute[0] == "197" {
				diskStats["currentPendingSectorCount"] = attribute[9]
			}

			if attribute[0] == "199" {
				diskStats["CRCErrorCount"] = attribute[9]
			}

			if attribute[0] == "233" {
				diskStats["mediaWearoutIndicator"] = attribute[3]
			}
		}
	}

	return diskStats, nil
}
