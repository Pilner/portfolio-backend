.PHONY: init vendor

init:
	@echo "Initializing local development environment..."
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit
	@echo "Git hooks configured successfully!"

run:
	@go build -o frv ./cmd/app
	@./frv

dockerup:
	docker-compose -f infra/docker-compose.yaml up -d
dockerdown:
	docker compose -f infra/docker-compose.yaml down -v --remove-orphans


db-fresh:
	docker exec -it postgres_db psql -U root -d postgres -c "DROP DATABASE IF EXISTS frv;"
	docker exec -it postgres_db psql -U root -d postgres -c "CREATE DATABASE frv;"
db-psql:
	psql -h localhost -U root -d frv

vendor:
	@echo "Vendoring dependencies..."
	go mod tidy
	go mod vendor
	@echo "Scrubbing vendor directory of build and test artifacts..."
	@find vendor -type f -name "Dockerfile*" -delete
	@find vendor -type f -name "docker-compose*" -delete
	@find vendor -type f -name "Makefile" -delete
	@find vendor -type f -name "*.md" -delete
	@find vendor -type f -name "*.txt" ! -name "modules.txt" -delete
	@echo "Vendor Cleanup Complete!"