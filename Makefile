PROJECT_NAME := "ci6ndex"
EXEC_NAME := civ

.PHONY: help ## print this
help:
	@echo ""
	@echo "$(PROJECT_NAME) Development CLI"
	@echo ""
	@echo "Usage:"
	@echo "  make <command>"
	@echo ""
	@echo "Commands:"
	@grep '^.PHONY: ' Makefile | sed 's/.PHONY: //' | awk '{split($$0,a," ## "); printf "  \033[34m%0-10s\033[0m %s\n", a[1], a[2]}'


.PHONY: update ## update dependencies
update:
	@echo "Updating dependencies..."
	@go get -u
	@go mod tidy
	@echo "Done!"

.PHONY: doctor ## checks if local environment is ready for development
doctor:
	@echo "Checking local environment..."
	@if ! command -v go &> /dev/null; then \
		echo "`go` is not installed. Please install it first."; \
		exit 1; \
	fi
	@if [[ ! ":$$PATH:" == *":$$HOME/go/bin:"* ]]; then \
		echo "GOPATH/bin is not in PATH. Please add it to your PATH variable."; \
		exit 1; \
	fi
	@if ! command -v sqlc &> /dev/null; then \
		echo "`sqlc` is not installed. Please run 'make deps'."; \
		exit 1; \
	fi

	@if ! command -v docker &> /dev/null; then \
		echo "`docker` is not installed. Please install it first."; \
		exit 1; \
	fi

	@echo "Local environment OK"


.PHONY: build ## build the project
build:
	$(MAKE) generate
	@echo "Building..."
	@go build -o ./bin/$(EXEC_NAME) .
	@echo "Done!"

.PHONY: clean ## delete generated code
clean:
	@echo "Deleting generated code..."
	@rm -rf ci6ndex/generated/
	@rm -rf bin/
	@echo "Done!"

.PHONY: generate ## generate server and database code
generate:
	@echo "Generating database models..."
	@sqlc generate -f sqlc.yaml
	@echo "Done!"

.PHONY: run ## run the project
run:
	$(MAKE) build
	@(export $$(cat .env | xargs) && ./bin/$(EXEC_NAME) bot serve)

.PHONY: sync ## sync discord commands with API
sync:
	$(MAKE) build
	@(export $$(cat .env | xargs) && ./bin/$(EXEC_NAME) bot sync)