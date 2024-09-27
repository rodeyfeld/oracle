# ORACLE
Feasibility and Archive Search Data. Written in Go.

To launch:
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