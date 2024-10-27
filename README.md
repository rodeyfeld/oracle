# ORACLE
Feasibility and Archive Search Generator Data. Written in Go.


Install Go:
https://go.dev/doc/install
`sudo snap install go --classic`

# Create .env file
```
DEBUG_MODE=true
DB_URL=mongodb://root.example@localhost:27017
#DOCKER_DB_URL=mongodb://root.example@mongo:27017
```


To launch:

# Local
```
go run .
```

### Docker
```
docker compose up --build
```


Currently serves four links:
- http://localhost:7777/attendPastCold: Scrapes the SENTINEL Database (Warning! Consumes massive resources until i figure out buffered channels)
- http://localhost:7777/attendPast: Returns a randomly generated SENTINEL object for testing purposes
- http://localhost:7777/attendPresent: Returns a randomly generated order object for testing purposes
- http://localhost:7777/attendFuture: Returns a randomly generated feasibility confidence score for testing purposes


Scraping the SENTINEL Database will load all data into a mongodb collection "catalogs". Quick access to this data can be found through mongo-express
- http://localhost:8888/