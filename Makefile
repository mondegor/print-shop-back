-include Makefile.mk

build:
	mrcmd install

build-api:
	mrcmd openapi build-all

deps:
	mrcmd go deps

deps-upgrade:
	mrcmd go get -u ./...
	mrcmd go tidy

migrate:
	mrcmd go-migrate up

generate:
	mrcmd go generate

lint:
	mrcmd go fmt
	mrcmd go fmti
	mrcmd go fmti2
	mrcmd go lint

test:
	mrcmd go test

test-report:
	mrcmd go test-report

plantuml:
	mrcmd plantuml build-all

app-conf:
	mrcmd config
	mrcmd docker-compose conf

app-start:
	mrcmd start

app-state:
	mrcmd docker-compose ps

app-logs:
	mrcmd docker-compose logs

app-stop:
	mrcmd stop

.PHONY: build build-api deps deps-upgrade migrate generate lint test test-report plantuml app-conf app-start app-state app-logs app-stop