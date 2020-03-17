.PHONY: build dist clean deps help
.DEFAULT_GOAL := help

# AutoDoc
define PRINT_HELP_PYSCRIPT
import re, sys

for line in sys.stdin:
	match = re.match(r'^([a-zA-Z_-]+):.*?## (.*)$$', line)
	if match:
		target, help = match.groups()
		print("%-10s %s" % (target, help))
endef
export PRINT_HELP_PYSCRIPT

build: ## Build the app
	@GOOS=linux go build -o ./build/main ./cmd/stm

dist: build ## Prepare the deployment package
	@mkdir -p ./dist
	@cd ./build
	@zip -q ./dist/function.zip ./main
	@cd ..

clean: ## Clean intermediate artifacts
	@rm -rf ./build ./dist

deps: ## Download the app dependencies
	@go mod vendor

help: ## Show this help
	@python -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)
