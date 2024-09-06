module oracle.com/audience

go 1.23.0

replace oracle.com/attendant => ./internal/attendant

replace oracle.com/soothsayer => ./internal/soothsayer

replace oracle.com/scholar => ./internal/scholar

replace oracle.com/chaos => ./internal/chaos

replace oracle.com/order => ./internal/order

require oracle.com/attendant v0.0.0-00010101000000-000000000000

require (
	oracle.com/chaos v0.0.0-00010101000000-000000000000 // indirect
	oracle.com/order v0.0.0-00010101000000-000000000000 // indirect
	oracle.com/scholar v0.0.0-00010101000000-000000000000 // indirect
	oracle.com/soothsayer v0.0.0-00010101000000-000000000000 // indirect
)
