package main

import (
	"fmt"
	"os"
	"strconv"

	zsend "github.com/blacked/go-zabbix"
	docopt "github.com/docopt/docopt-go"
)

const (
	version  = "[manual build]"
	sudo     = "/usr/bin/sudo"
	smartctl = "/usr/bin/smartctl"
	SAS      = "SAS"
	SATA     = "SATA"
)

func main() {
	usage := `zabbix-agent-extension-smartctl

Usage:
	zabbix-agent-extension-smartctl [options]

Options:
	-d --disk <name>          Disk name [default: /dev/sda].
	-i --interface <type>	  Interface type [default: SATA].
	-z --zabbix-host <zhost>  Hostname or IP address of zabbix server 
							   [default: 127.0.0.1].
	-p --zabbix-port <zport>  Port of zabbix server [default: 10051].
	--zabbix-prefix <prefix>  Add part of your prefix for key [default: None].
	--discovery               Run low-level discovery for determine disks.
	-h --help                 Show this screen.
`
	args, _ := docopt.Parse(usage, nil, true, version, false)

	if args["--discovery"].(bool) {
		if err := discovery(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	zabbixHost := args["--zabbix-host"].(string)
	zabbixPort, err := strconv.Atoi(args["--zabbix-port"].(string))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	zabbixPrefix := args["--zabbix-prefix"].(string)
	if zabbixPrefix == "None" {
		zabbixPrefix = "SMART.disk"
	} else {
		zabbixPrefix = fmt.Sprintf("%s.%s", zabbixPrefix, "SMART.disk")
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	diskInfo := make(map[string]string)
	diskInfo["name"] = args["--disk"].(string)
	diskInfo["interface"] = args["--interface"].(string)

	diskStats, err := getDiskStats(diskInfo)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var zabbixMetrics []*zsend.Metric
	zabbixMetrics = createMetrics(diskStats, hostname, zabbixPrefix, zabbixMetrics, diskInfo["name"])

	packet := zsend.NewPacket(zabbixMetrics)
	sender := zsend.NewSender(
		zabbixHost,
		zabbixPort,
	)
	sender.Send(packet)

	fmt.Println("OK")
}
