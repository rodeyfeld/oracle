# ORACLE
Feasibility and Archive Search Generator Data. Written in Go.


Install Go:
https://go.dev/doc/install
`sudo snap install go --classic`

# Create .env file
```
DEBUG_MODE=false
SCRAPE_MODE=false
MONGO_DB_URL=mongodb://root:example@localhost:27017/
POSTGRES_DB_URL=postgresql://postgres:mypassyword@localhost:5432/lore
```


To launch:

# Local
```
go run .
```
If `SCRAPE_MODE=true` the scraper will attempt to get all results before starting the server


### Docker
```
docker compose up --build
```


Currently serves three links:
- http://localhost:7777/attendPast: Returns a randomly generated SENTINEL object for testing purposes
- http://localhost:7777/attendPresent: Returns a randomly generated order object for testing purposes
- http://localhost:7777/attendFuture: Returns a randomly generated feasibility confidence score for testing purposes


Scraping the SENTINEL Database will load all data into a mongodb collection "catalogs". Quick access to this data can be found through mongo-express
- http://localhost:8888/