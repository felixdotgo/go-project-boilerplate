.DEFAULT_GOAL = help

.PHONY: docker.down
docker.down: ## Run docker-compose down
	@docker-compose down


.PHONY: docker.up
docker.up: ## Run docker-compose up
	@docker-compose up


.PHONY: docker.upd
docker.upd: ## Run docker-compose up -d
	@docker-compose up -d


.PHONY: migrate.down
migrate.down: ## Revert all down migrations
	@go run cmd/migrator/migrator.go down


.PHONY: migrate.up
migrate.up: ## Run all migrations
	@go run cmd/migrator/migrator.go up


.PHONY: vendor
vendor: ## Run go mod vendor
	go mod vendor


.PHONY: help
help:  ## Display this help
	@echo "Usage: make \033[36m<command>\033[0m"
	@awk 'BEGIN {FS = ":.*##"; printf "\nCommands:\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
