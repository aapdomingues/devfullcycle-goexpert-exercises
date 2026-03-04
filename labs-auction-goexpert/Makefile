COMPOSE_TEST=docker compose -f docker-compose.test.yml

test: up-test wait-mongo go-test down-test

up-test:
	$(COMPOSE_TEST) up -d

down-test:
	$(COMPOSE_TEST) down -v

go-test:
	go test ./internal/infra/database/auction -run TestAuctionShouldCloseAutomatically -count=1

wait-mongo:
	@echo "Waiting for MongoDB..."
	@until $(COMPOSE_TEST) exec -T mongodb mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "MongoDB ready!"