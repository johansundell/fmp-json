GHACCOUNT := johansundell
NAME := fmp-json
VERSION := v0.6

include common.mk

deps:
	go get github.com/c4milo/github-release
	go get github.com/mitchellh/gox


