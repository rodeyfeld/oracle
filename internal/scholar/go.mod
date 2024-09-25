module oracle.com.internal/scholar

go 1.23.0

replace oracle.com/order => ../order

replace oracle.com/chaos => ../chaos

replace oracle.com/copernicus => ./copernicus

require (
	oracle.com/chaos v0.0.0-00010101000000-000000000000
	oracle.com/copernicus v0.0.0-00010101000000-000000000000
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
	oracle.com/order v0.0.0-00010101000000-000000000000 // indirect
)
