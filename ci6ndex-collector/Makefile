PROJECT_NAME := "ci6ndex-collector"
EXEC_NAME := ci6ndex-collector

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


.PHONY: doctor ## checks if local environment is ready for development
doctor:
	@echo "Checking local environment..."
	@if ! command -v wails &> /dev/null; then \
		echo "`wails` CLI is not installed. Please install it first."; \
		exit 1; \
	fi
	@if [[ ! ":$$PATH:" == *":$$HOME/go/bin:"* ]]; then \
		echo "GOPATH/bin is not in PATH. Please add it to your PATH variable."; \
		exit 1; \
	fi
	@wails doctor

	@echo "Local environment OK"


.PHONY: build ## build the project
build:
	@echo "Building..."
	@wails build
	@echo "Done!"

.PHONY: start ## run the project
start:
	@echo "Running..."
	@wails dev
	@echo "Done!"