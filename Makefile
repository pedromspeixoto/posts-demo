# General purpose targets.

.PHONY: help
help: ## Display available commands.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker targets.

.PHONY: docker-compose-up-local
docker-compose-up-local: ## Run all services locally using docker compose
	docker-compose -f docker-compose.local.yaml up -d --build

.PHONY: docker-compose-down-local
docker-compose-down-local: ## Delete all services locally using docker compose
	docker-compose -f docker-compose.local.yaml down