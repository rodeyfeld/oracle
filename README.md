# Oracle

Scrapes Sentinel satellite archives and dumps metadata into databases.

## Setup

Get Go:
```bash
sudo snap install go --classic
```

Make a `.env`:
```bash
DEBUG_MODE=false
SCRAPE_MODE=false
MONGO_DB_URL=mongodb://root:example@localhost:27017/
POSTGRES_DB_URL=postgresql://postgres:mypassword@localhost:5432/augur
```

## Run it

```bash
go run cmd/main.go 
```

Scrapes Sentinel archive APIs, stores raw metadata in MongoDB, syncs processed data to Postgres so Augur can query it. Basically builds a local cache so we're not constantly hitting external APIs and running into rate limits.
