.PHONY: all clean-all build cleand-deps deps ver

DATE := $(shell git log -1 --format="%cd" --date=short | sed s/-//g)
COUNT := $(shell git rev-list --count HEAD)
COMMIT := $(shell git rev-parse --short HEAD)

BINARYNAME := zabbix-agent-extension-smartctl
CONFIG := zabbix-agent-extension-smartctl.conf
SUDOERS := zabbix-smartctl
VERSION := "${DATE}.${COUNT}_${COMMIT}"

LDFLAGS := "-X main.version=${VERSION}"

default: all

all: clean-all deps build 

ver:
	@echo ${VERSION}

clean-all: clean-deps
	@echo Clean builded binaries
	rm -rf .out/
	@echo Done

build:
	@echo Build
	ln -s ${PWD}/vendor/ ${PWD}/vendor/src
	GOPATH="${PWD}/vendor" go build -v -o .out/${BINARYNAME} -ldflags ${LDFLAGS} *.go
	@echo Done

clean-deps:
	@echo Clean dependencies
	rm -rf vendor/*

deps:
	@echo Fetch dependencies
	git submodule update --init

install:
	@echo Install
	cp .out/${BINARYNAME} /usr/bin/${BINARYNAME}
	cp ${CONFIG} /etc/zabbix/zabbix_agentd.conf.d/${CONFIG}
	cp ${SUDOERS} /etc/sudoers.d/${CONFIG}
	@echo Done 

uninstall:
	@echo Remove
	rm /usr/bin/${BINARYNAME}
	rm /etc/zabbix/zabbix_agentd.conf.d/${CONFIG}
	rm /etc/sudoers.d/${SUDOERS}
	@echo Done
