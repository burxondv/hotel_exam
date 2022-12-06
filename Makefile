POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=hotel_exam

-include .env
  
DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"


print:
	echo "$(DB_URL)"
	
swag-init:
	swag init -g api/api.go -o api/docs

start:
	go run main/main.go

migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1


local-up:
	docker compose --env-file ./.env.docker up -d

.PHONY: start migrateup migratedown