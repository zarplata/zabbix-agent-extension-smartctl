# zabbix-agent-extension-smartctl

zabbix-agent-extension-smartctl - this zabbix agent extension for monitoring disk with smartctl.

### Supported features

This extension supports the detection of disks on a Linux host and getting the statistics from disk.

### Installation

#### Manual build

```sh
# Building
git clone https://github.com/zarplata/zabbix-agent-extension-smartctl.git
cd zabbix-agent-extension-smartctl
make

#Installing
make install

# By default, binary installs into /usr/bin/ and zabbix config in /etc/zabbix/zabbix_agentd.conf.d/ but,
# you may manually copy binary to your executable path and zabbix config to specific include directory
```

#### Arch Linux package
```sh
# Building
git clone https://github.com/zarplata/zabbix-agent-extension-smartctl.git
git checkout pkgbuild

makepkg

#Installing
pacman -U *.tar.xz
```

### Dependencies

zabbix-agent-extension-smartcl requires [zabbix-agent](http://www.zabbix.com/download) v2.4+ to run.

### Zabbix configuration
In order to start getting metrics, it is enough to import template and attach it to monitored node.

`WARNING:` You must define macro with name - `{$ZABBIX_SERVER_IP}` in global or local (template) scope with IP address of  zabbix server.
