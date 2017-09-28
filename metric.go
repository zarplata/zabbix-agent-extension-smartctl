package main

import (
	"fmt"
	"strings"

	zsend "github.com/blacked/go-zabbix"
)

func getPrefix(zabbixPrefix, key, diskName string) string {
	if strings.Contains(diskName, ",") {
		return fmt.Sprintf("%s.%s.[\"%s\"]", zabbixPrefix, key, diskName)
	}

	return fmt.Sprintf("%s.%s.[%s]", zabbixPrefix, key, diskName)
}

func createMetrics(
	stats map[string]string,
	hostname, zabbixPrefix string,
	zabbixMetrics []*zsend.Metric,
	diskName string,
) []*zsend.Metric {
	for statName, stat := range stats {
		zabbixMetrics = append(
			zabbixMetrics,
			zsend.NewMetric(
				hostname,
				getPrefix(zabbixPrefix, statName, diskName),
				stat,
			),
		)
	}

	return zabbixMetrics
}
