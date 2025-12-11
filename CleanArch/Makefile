.PHONY: createmigration migrate migratedown up down wait-mysql migrate run init

APP_DIR=cmd/ordersystem

# 1) Sobe MySQL e RabbitMQ em segundo plano
up:
	docker-compose up -d

# (Opcional) Derruba os containers
down:
	docker-compose down

# 2) Espera o MySQL estar pronto antes de rodar as migrations
wait-mysql:
	@echo "⏳ Aguardando MySQL ficar pronto..."
	@until docker-compose exec -T mysql mysqladmin ping -h "localhost" --silent; do \
		echo "."; \
		sleep 2; \
	done
	@echo "✅ MySQL está pronto!"

wait-rabbitmq:
	@echo "⏳ Aguardando RabbitMQ ficar pronto..."
	@until docker-compose exec -T rabbitmq rabbitmq-diagnostics -q ping > /dev/null 2>&1; do \
		echo "."; \
		sleep 2; \
	done
	@echo "✅ RabbitMQ está pronto!"

# 3) Roda migrations usando o binário 'migrate' instalado na sua máquina
migrate: up wait-mysql wait-rabbitmq
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

# 4) Roda sua aplicação Go (fora do Docker)
run:
	go run main.go

# 5) Target completo: sobe tudo, migra e roda a app
init: migrate
	cd $(APP_DIR) && go run main.go wire_gen.go


createmigration:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

migratedown:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down



