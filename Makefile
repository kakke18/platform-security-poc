.PHONY: help proto dev up down logs clean

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

proto: ## Generate code from Protocol Buffers
	buf generate

dev: ## Run project-api locally (requires Go 1.23+)
	cd project-api && go run cmd/server/main.go

up: ## Start services with Docker Compose
	docker-compose up --build

down: ## Stop services
	docker-compose down

logs: ## Show logs
	docker-compose logs -f

clean: ## Clean generated files and containers
	docker-compose down -v
	find . -name "*.pb.go" -delete
	find . -name "*connect.go" -delete

install-deps: ## Install dependencies (project-api)
	cd project-api && go mod download
