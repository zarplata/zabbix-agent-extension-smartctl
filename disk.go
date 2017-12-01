package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	hierr "github.com/reconquest/hierr-go"
)

func hasBit(n int64, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func errorFromBitMask(diskName string, err error) error {

	var exitCodeSmart = []string{
		"Bit 0: Command line did not parse",
		`Bit 1: Device open failed, device did not return an IDENTIFY
DEVICE structure, or device is in a low-power mode (see '-n' option above)`,
		`Bit 2: Some SMART command to the disk failed, or there was a
checksum error in a SMART data structure (see '-b' option above)`,
		"Bit 3: SMART status check returned DISK FAILING",
		"Bit 4: We found prefail Attributes <= threshold",
		`Bit 5: SMART status check returned DISK OK but we found that
some (usage or prefail) Attributes have been <= threshold
at some time in the past`,
		"Bit 6: The device error log contains records of errors",
		`Bit 7: The device self-test log contains records of errors.
[ATA only] Failed self-tests outdated by a newer successful extended
self-test are ignored`,
	}

	var errorSmart string
	var waitStatus syscall.WaitStatus

	if exitError, ok := err.(*exec.ExitError); ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
		bitMaskInt64 := int64(waitStatus.ExitStatus())

		for i := uint(0); i < uint(len(exitCodeSmart)); i++ {
			if hasBit(bitMaskInt64, i) {
				errorSmart = fmt.Sprintf("%s%s\n",
					errorSmart,
					exitCodeSmart[i],
				)
				if i <= 1 {
					hierr.Fatalf(
						fmt.Errorf(strings.TrimRight(errorSmart, "\n")),
						"can't get stats for %s disk",
						diskName,
					)
				}
			}
		}
		return hierr.Errorf(
			fmt.Errorf(strings.TrimRight(errorSmart, "\n")),
			"smartctl return not null exit code for %s disk",
			diskName,
		)
	}

	return hierr.Errorf(
		fmt.Errorf("can't get abnormal exit code"),
		"smartctl return not null exit code for %s disk",
		diskName,
	)
}

func getDiskStats(diskInfo map[string]string) (map[string]string, error) {

	args := []string{smartctl, "-a"}
	args = append(args, strings.Split(diskInfo["name"], " ")...)
	out, err := exec.Command(sudo, args...).CombinedOutput()

	var exitMsg error

	if err != nil {
		exitMsg = errorFromBitMask(diskInfo["name"], err)
	}

	if diskInfo["interface"] == SAS {
		return getSASStats(diskInfo, string(out)), exitMsg
	}
	return getSATAStats(diskInfo, string(out)), exitMsg

}

func getSASStats(
	diskInfo map[string]string,
	rawStats string,
) map[string]string {
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

	return diskStats
}

func getSATAStats(
	diskInfo map[string]string,
	rawStats string,
) map[string]string {
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

	return diskStats
}
