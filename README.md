# ORACLE
SENTINEL 1/2 Satellite Archive Imagery Scraper. 


## Install Go:
https://go.dev/doc/install
`sudo snap install go --classic`

## Create .env file
```
DEBUG_MODE=false
SCRAPE_MODE=false
MONGO_DB_URL=mongodb://root:example@localhost:27017/
POSTGRES_DB_URL=postgresql://postgres:mypassyword@localhost:5432/augur
```

To launch:

## Local
```
go run cmd/main.go 
```

