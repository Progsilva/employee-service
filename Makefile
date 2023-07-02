.PHONY: test
gen: ## Run Tests into the packages
	@echo "Generate graphql schema"
	go run github.com/99designs/gqlgen generate

test: ## Run Tests into the packages
	@echo "Running tests"
	go test ./...

run: ## Run Tests into the packages
	@echo "run project"
	docker-compose up -d --build

stop: ## Run Tests into the packages
	@echo "stop project"
	docker-compose down