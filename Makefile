# Define standard colors
YELLOW := $(shell tput setaf 3)
GREEN := $(shell tput setaf 2)
BLUE := $(shell tput setaf 4)
RED := $(shell tput setaf 1)
MAGENTA := $(shell tput setaf 5)
CYAN := $(shell tput setaf 6)
WHITE := $(shell tput setaf 7)
RESET := $(shell tput sgr0)
BOLD := $(shell tput bold)

.DEFAULT_GOAL = help

OS := $(shell uname -s)

####################################################################################################
# Utilities
####################################################################################################
.PHONY: check_cmd_var
check_cmd_var:
	@if [ -z "$(CMD)" ]; then \
		echo "❌ $(RED)Error: CMD is not set$(RESET)"; \
		exit 1; \
	fi

.PHONY: check_docker_file
check_docker_file:
	@if [ ! -f "cmd/$(CMD)/Dockerfile" ]; then \
		echo "❌ $(RED)Error: Dockerfile doesn't exist in cmd/$(CMD). Cannot build.$(RESET)"; \
		exit 1; \
	fi


.PHONY: help
help:  ## Display this help
	@echo "$(GREEN)$(BOLD)"
	@cat scripts/banner.txt
	@echo "$(RESET)"
	@echo "$(WHITE)Usage:$(RESET) make $(CYAN)<command>$(RESET)"
	@echo "$(WHITE)Commands:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"; } \
		/^[.a-zA-Z_-]+:.*?##/ { printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(WHITE)%s:$(RESET)\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@echo "\n$(WHITE)For detailed information about a command, add $(YELLOW)--help$(WHITE) after the command.$(RESET)"


####################################################################################################
# Commands to run core services for development
####################################################################################################
.PHONY: docker.down
docker.down: ## Run docker-compose down
	@echo "🚧 $(BLUE)Running docker-compose down...$(RESET)"
	@cd devenv && docker-compose down
	@echo "🎉 $(GREEN)Docker containers stopped and removed$(RESET)"


.PHONY: docker.up
docker.up: ## Run docker-compose up -d
	@echo "🚀 $(BLUE)Starting docker-compose in detached mode...$(RESET)"
	@cd devenv && docker-compose up -d
	@echo "🎉 $(GREEN)Docker containers started in background$(RESET)"


####################################################################################################
# Development commands
####################################################################################################
.PHONY: certs
certs: ## Generate SSL certificates with mkcert for .loc domains
	@echo "$(CYAN)Creating certificates directory...$(RESET)"
	@mkdir -p devenv/certs
	@echo "$(CYAN)Generating SSL certificates for goproject.local domains...$(RESET)"
	@cd devenv/certs && mkcert -cert-file goproject.local.pem -key-file goproject.local-key.pem "*.goproject.local" goproject.local traefik.goproject.local
	@echo "$(CYAN)Installing mkcert CA...$(RESET)"
	@mkcert -install
	@echo "$(GREEN)SSL certificates generated successfully!$(RESET)"
	@echo "$(YELLOW)Certificate files:$(RESET)"
	@echo "  - $(WHITE)certs/goproject.local.pem$(RESET) (certificate)"
	@echo "  - $(WHITE)certs/goproject.local-key.pem$(RESET) (private key)"



.PHONY: vendor
vendor: ## Run go mod vendor
	@echo "✨ $(BLUE)Running go mod vendor...$(RESET)"
	go mod vendor
	@echo "🎉 $(GREEN)Vendor directory updated$(RESET)"


.PHONY: run
run: check_cmd_var ## Run specific dir inside `cmd` with `make run CMD=<your dir>`
	@echo "✨ $(BLUE)Running cmd/$(CMD)...$(RESET)"
	@go run "cmd/$(CMD)/main.go"


.PHONY: generate-proto
generate-proto: ## Generate protobuf code
	@echo "✨ $(BLUE)Generating protobuf code...$(RESET)"
	@buf generate --include-imports
	@echo "🎉 $(GREEN)Protobuf code generated successfully$(RESET)"


.PHONY: build
build: check_cmd_var check_docker_file ## Build specific dir inside `cmd` with `make build CMD=<your dir> [PLATFORM=linux/amd64] [BUILD_ARGS="KEY1=value1,KEY2=value2"]`
	@echo "✨ $(BLUE)Building $(CMD) Docker image...$(RESET)"
	$(eval TIMESTAMP := $(shell date +%Y%m%d-%H%M%S))
	$(eval PLATFORM ?= linux/amd64)
	$(eval PARSED_ARGS := $(shell echo "$(BUILD_ARGS)" | sed 's/,/ --build-arg /g'))
	$(eval BUILD_ARGUMENTS := $(if $(BUILD_ARGS), $(PARSED_ARGS),))

	@cd cmd/$(CMD)/ && docker-compose build --build-arg PLATFORM="$(PLATFORM)" $(BUILD_ARGUMENTS) || { echo "❌ $(RED)Docker build failed for $(CMD)$(RESET)" && cd - > /dev/null; exit 1; }
	@echo "🎉 $(GREEN)Successfully built $(CMD) image$(RESET)"


.PHONY: up
up: build ## up: Build and start the service
	@echo "🚀 $(BLUE)Starting $(CMD) service...$(RESET)"
	@cd cmd/$(CMD)/ && docker-compose up


.PHONY: upd
upd: build ## up: Build and start the service (detached mode)
	@echo "🚀 $(BLUE)Starting $(CMD) service...$(RESET)"
	@cd cmd/$(CMD)/ && docker-compose up -d


.PHONY: install
install: ## Install dependencies
	@echo "✨ $(BLUE)Installing dependencies...$(RESET)"
	@./scripts/setup_environment.sh
	@echo "🎉 $(GREEN)Dependencies installed successfully$(RESET)"

