-include Makefile.mk

.PHONY: build
build:
	mrcmd install

.PHONY: build-api
build-api:
	mrcmd openapi build-all

.PHONY: deps
deps:
	mrcmd go deps

.PHONY: deps-upgrade
deps-upgrade:
	mrcmd go get -u ./...
	mrcmd go tidy

.PHONY: migrate
migrate:
	mrcmd go-migrate up

.PHONY: generate
generate:
	mrcmd go generate

.PHONY: fmt
fmt:
	mrcmd go fmt

.PHONY: fmti
fmti:
	mrcmd go fmti

.PHONY: lint
lint:
	mrcmd golangci-lint check

.PHONY: test
test:
	mrcmd go test

.PHONY: test-report
test-report:
	mrcmd go test-report

.PHONY: plantuml
plantuml:
	mrcmd plantuml build-all

.PHONY: app-conf
app-conf:
	mrcmd config
	mrcmd docker-compose conf

.PHONY: app-start
app-start:
	mrcmd start

.PHONY: app-state
app-state:
	mrcmd docker-compose ps

.PHONY: app-logs
app-logs:
	mrcmd docker-compose logs

.PHONY: app-stop
app-stop:
	mrcmd stop