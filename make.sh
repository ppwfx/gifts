#!/usr/bin/env bash

set -eo pipefail

function run {
	EMPLOYEES_JSON=employees.json GIFTS_JSON=gifts.json ADDR=:30030 go run main.go
}

function test {
	go test -v
}