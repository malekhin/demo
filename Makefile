.PHONY: build

PROJECT_NAME=demo
COMPOSE_DEV=./build/dev/docker-compose.yaml

DB_NAME=main
BASE_BRANCH ?= main
FAIL_UNDER ?= 60

build:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) up --build -d

up:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) up -d

down:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) down

restart:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) restart

swag:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec app sh -c "swag init --parseDependency -g cmd/app/main.go"

test:
	docker compose -p ${PROJECT_NAME} -f ${COMPOSE_DEV} exec app sh -c "ENVIRONMENT=test go test -v -p=1 -count=1 ./..."

test-coverage:
	docker compose -p ${PROJECT_NAME} -f ${COMPOSE_DEV} exec app sh -c "ENVIRONMENT=test go test -cover -coverprofile ./tmp/coverage.out ./..."
	docker compose -p ${PROJECT_NAME} -f ${COMPOSE_DEV} exec app sh -c "cat ./tmp/coverage.out | grep -v -e "/mocks" -e "/tests" > ./tmp/coverage.out.filtered"
	docker compose -p ${PROJECT_NAME} -f ${COMPOSE_DEV} exec app sh -c "gocover-cobertura < ./tmp/coverage.out.filtered > ./tmp/coverage.xml"
	docker compose -p ${PROJECT_NAME} -f ${COMPOSE_DEV} exec app sh -c "diff-cover ./tmp/coverage.xml --compare-branch=origin/$(BASE_BRANCH) --fail-under=$(FAIL_UNDER)"

test-visualization:
	docker compose -p ${PROJECT_NAME} -f ${COMPOSE_DEV} exec app sh -c "go tool cover -html=./coverage.out -o coverage.out.html"

shell:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec app bash

redis:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec redis redis-cli

db:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec postgres psql $(DB_NAME)

logs:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) logs -f --tail 100


generate:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec app sh -c "go generate ./..."

lint:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec app sh -c "golangci-lint run"

drop:
	docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec app /migrate_pg.sh drop

migration:
	@read -p "Migration name: " migration; \
		docker compose -p ${PROJECT_NAME} -f $(COMPOSE_DEV) exec app sh -c "/migrate_pg.sh create -ext sql -dir /app/migrations $$migration"
