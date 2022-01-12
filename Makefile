network:
	docker network create go-song-network

postgres:
	docker run --name postgres12 --network go-song-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root go_song

dropdb:
	docker exec -it postgres12 dropdb go_song

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/go_song?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/go_song?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/go_song?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/go_song?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate
