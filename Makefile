version := 0.0.2
.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Makefile Commands:"
	@echo "----------------------------------------------------------------"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo "----------------------------------------------------------------"

.PHONY: up
up: ## start docker containers
	@docker-compose -f docker-compose.yml pull
	@docker-compose -f docker-compose.yml up -d

.PHONY: down
down: ## shuts down docker containers
	docker-compose -f docker-compose.yml down --remove-orphans

run: ## run server
	@go run main.go

version: ## iterate sem-ver
	go generate
	bumpversion patch --allow-dirty

tag: ## tag sem-ver
	git tag v$(version)

push: ## push updated code to github
	git push origin master
	git push origin v$(version)

test: ## run tests
	@go test -v