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
	oracle.com/chaos v0.0.0-00010101000000-000000000000 // indirect
	oracle.com/copernicus v0.0.0-00010101000000-000000000000 // indirect
	oracle.com/order v0.0.0-00010101000000-000000000000 // indirect
)
