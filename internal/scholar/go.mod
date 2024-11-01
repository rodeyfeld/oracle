module github.com/rodeyfeld/oracle/internal/scholar

go 1.23.0

replace github.com/rodeyfeld/oracle/order => ../order

replace github.com/rodeyfeld/oracle/chaos => ../chaos

replace github.com/rodeyfeld/oracle/copernicus => ./copernicus

require (
	github.com/rodeyfeld/oracle/chaos v0.0.0-00010101000000-000000000000
	github.com/rodeyfeld/oracle/copernicus v0.0.0-00010101000000-000000000000
	github.com/rodeyfeld/oracle/order v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver/v2 v2.0.0-beta2
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.1 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/venicegeo/geojson-go v0.1.0 // indirect
	github.com/venicegeo/pzsvc-lib v0.0.0-20161208182529-fca89502ff2c // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)
