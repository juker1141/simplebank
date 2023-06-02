postgres:
	docker-compose up -d

createdb:
	docker-compose exec -it db createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker-compose exec -it db dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/juker1141/simplebank/db/sqlc store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock