GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=greenhouse_backend
VERSION?=0.0.0

SRV_DEPLOY=192.168.0.27

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build

all: help

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## Build:
build: build-server ## Build the project and put the output binary in out/bin/
	@mkdir -p out
	@cp greenhouse.toml out/

build-server: ## Build only the golang server
	@mkdir -p out
	$(GOCMD) build -v -o out/$(BINARY_NAME) .
	@cp greenhouse.service out/

clean: ## Remove build related file
	rm -fr ./bin
	rm -fr ./out
	rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

## Upload:
upload: upload-server ## Upload binaries using ssh

upload-server: build-server ## Upload server binary
	@rsync -avP ./out/$(BINARY_NAME) root@$(SRV_DEPLOY):/usr/local/bin/$(BINARY_NAME)
	@ssh root@$(SRV_DEPLOY) chown 755 /usr/local/bin/$(BINARY_NAME)

## Test:
test: ## Run the tests of the project
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov