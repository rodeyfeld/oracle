module oracle.com/audience

go 1.23.0

replace oracle.com/soothsayer => ./internal/soothsayer

replace oracle.com/scholar => ./internal/scholar

replace oracle.com/copernicus => ./internal/scholar/copernicus

replace oracle.com/chaos => ./internal/chaos

replace oracle.com/order => ./internal/order

replace oracle.com/bazaar => ./internal/bazaar

require (
	oracle.com/bazaar v0.0.0-00010101000000-000000000000
	oracle.com/scholar v0.0.0-00010101000000-000000000000
	oracle.com/soothsayer v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver/v2 v2.0.0-beta2 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	oracle.com/chaos v0.0.0-00010101000000-000000000000 // indirect
	oracle.com/copernicus v0.0.0-00010101000000-000000000000 // indirect
	oracle.com/order v0.0.0-00010101000000-000000000000 // indirect
)
