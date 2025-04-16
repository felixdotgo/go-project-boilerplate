# Define standard colors
YELLOW := $(shell tput setaf 3)
GREEN := $(shell tput setaf 2)
BLUE := $(shell tput setaf 4)
RED := $(shell tput setaf 1)
MAGENTA := $(shell tput setaf 5)
CYAN := $(shell tput setaf 6)
WHITE := $(shell tput setaf 7)
RESET := $(shell tput sgr0)

.DEFAULT_GOAL = help

.PHONY: docker.down
docker.down: ## Run docker-compose down
	@echo "🚧 $(BLUE)Running docker-compose down...$(RESET)"
	@docker-compose down
	@echo "🎉 $(GREEN)Docker containers stopped and removed$(RESET)"


.PHONY: docker.up
docker.up: ## Run docker-compose up
	@echo "🚀 $(BLUE)Starting docker-compose...$(RESET)"
	@docker-compose up


.PHONY: docker.upd
docker.upd: ## Run docker-compose up -d
	@echo "🚀 $(BLUE)Starting docker-compose in detached mode...$(RESET)"
	@docker-compose up -d
	@echo "🎉 $(GREEN)Docker containers started in background$(RESET)"


.PHONY: migrate.down
migrate.down: ## Revert all down migrations
	@echo "✨ $(YELLOW)Reverting all migrations...$(RESET)"
	@go run cmd/migrator/migrator.go down
	@echo "🎉 $(GREEN)Migrations reverted successfully$(RESET)"


.PHONY: migrate.up
migrate.up: ## Run all migrations
	@echo "🚀 $(YELLOW)Running all migrations...$(RESET)"
	@go run cmd/migrator/migrator.go up
	@echo "🎉 $(GREEN)Migrations completed successfully$(RESET)"


.PHONY: vendor
vendor: ## Run go mod vendor
	@echo "✨ $(BLUE)Running go mod vendor...$(RESET)"
	go mod vendor
	@echo "🎉 $(GREEN)Vendor directory updated$(RESET)"


.PHONY: cmdcheck
cmdcheck:
	@if [ -z "$(CMD)" ]; then \
		echo "❌ $(RED)Error: CMD is not set$(RESET)"; \
		exit 1; \
	fi


.PHONY: run
run: cmdcheck ## Run specific dir inside `cmd` with `make run CMD=<your dir>`
	@echo "✨ $(BLUE)Running cmd/$(CMD)...$(RESET)"
	@go run "cmd/$(CMD)/main.go"


.PHONY: dockerfile_check
dockerfile_check:
	@if [ ! -f "cmd/$(CMD)/Dockerfile" ]; then \
		echo "❌ $(RED)Error: Dockerfile doesn't exist in cmd/$(CMD). Cannot build.$(RESET)"; \
		exit 1; \
	fi


.PHONY: generate-proto
generate-proto: ## Generate protobuf code
	@echo "✨ $(BLUE)Generating protobuf code...$(RESET)"
	@buf generate --include-imports
	@echo "🎉 $(GREEN)Protobuf code generated successfully$(RESET)"


.PHONY: build
build: cmdcheck dockerfile_check ## Build specific dir inside `cmd` with `make build CMD=<your dir> [PLATFORM=linux/amd64] [BUILD_ARGS="KEY1=value1,KEY2=value2"]`
	@echo "✨ $(BLUE)Building $(CMD) Docker image...$(RESET)"
	$(eval TIMESTAMP := $(shell date +%Y%m%d-%H%M%S))
	$(eval PLATFORM ?= linux/amd64)
	$(eval PARSED_ARGS := $(shell echo "$(BUILD_ARGS)" | sed 's/,/ --build-arg /g'))
	$(eval BUILD_ARGUMENTS := $(if $(BUILD_ARGS), $(PARSED_ARGS),))

	@cd cmd/$(CMD)/ && docker-compose build --build-arg PLATFORM="$(PLATFORM)" $(BUILD_ARGUMENTS) || { echo "❌ $(RED)Docker build failed for $(CMD)$(RESET)" && cd - > /dev/null; exit 1; }
	@echo "🎉 $(GREEN)Successfully built $(CMD) image$(RESET)"


.PHONY: help
help:  ## Display this help
	@echo "$(GREEN)╔══════════════════════════════════════════════════════════════╗$(RESET)"
	@echo "$(GREEN)║                  $(WHITE)Go Project Boilerplate$(RESET)                      $(GREEN)║$(RESET)"
	@echo "$(GREEN)╚══════════════════════════════════════════════════════════════╝$(RESET)"
	@echo "$(WHITE)Usage:$(RESET) make $(CYAN)<command>$(RESET)"
	@echo "$(WHITE)Commands:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"; } \
		/^[.a-zA-Z_-]+:.*?##/ { printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(WHITE)%s:$(RESET)\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@echo "\n$(WHITE)For detailed information about a command, add $(YELLOW)--help$(WHITE) after the command.$(RESET)"
