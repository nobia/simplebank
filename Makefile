postgres:
	podman run --name postgres1 -d -p 5432:5432 -e POSTGRES_USER=u1 -e POSTGRES_PASSWORD=pass1 docker.io/library/postgres:14.5-alpine
createdb:
	podman exec -- postgres1 createdb --username=u1 --owner=u1 simple_bank
dropdb:
	podman exec -- postgres1 dropdb --username=u1 simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://u1:pass1@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://u1:pass1@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mockgen:
	mockgen -packge mockdb -destination db/mock/store.go github.com/nobia/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mockgen
