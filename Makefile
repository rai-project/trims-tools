all: generate

fmt: ## Formats code
	go fmt ./...

install-deps: ## Install dependencies
	go get golang.org/x/tools/cmd/goimports
	go get -u github.com/knq/xo

glide-install: ## Performs glide install
	glide install -v --force

logrus-fix: ## Fixes logrus
	rm -fr vendor/github.com/Sirupsen
	find vendor -type f -exec sed -i 's/Sirupsen/sirupsen/g' {} +

generate-models: ## Generates Models for NVPRof
	rm -fr nvprof/models
	xo 'file:./_fixtures/profile.timeline.nvprof?loc=auto' -o nvprof/models

generate: generate-models

help: ## Shows this help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: help

.DEFAULT_GOAL := generate
