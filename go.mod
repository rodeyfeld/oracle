module github.com/rodeyfeld/oracle/audience

go 1.23.0

replace github.com/rodeyfeld/oracle/soothsayer => ./internal/soothsayer

replace github.com/rodeyfeld/oracle/scholar => ./internal/scholar

replace github.com/rodeyfeld/oracle/copernicus => ./internal/scholar/copernicus

replace github.com/rodeyfeld/oracle/chaos => ./internal/chaos

replace github.com/rodeyfeld/oracle/order => ./internal/order

replace github.com/rodeyfeld/oracle/bazaar => ./internal/bazaar

require (
	github.com/joho/godotenv v1.5.1
	github.com/rodeyfeld/oracle/bazaar v0.0.0-00010101000000-000000000000
	github.com/rodeyfeld/oracle/scholar v0.0.0-00010101000000-000000000000
	github.com/rodeyfeld/oracle/soothsayer v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.1 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/rodeyfeld/oracle/chaos v0.0.0-00010101000000-000000000000 // indirect
	github.com/rodeyfeld/oracle/copernicus v0.0.0-00010101000000-000000000000 // indirect
	github.com/rodeyfeld/oracle/order v0.0.0-00010101000000-000000000000 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.11.4 // indirect
	go.mongodb.org/mongo-driver/v2 v2.0.0-beta2 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)
