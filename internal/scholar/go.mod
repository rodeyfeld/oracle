module oracle.com.internal/scholar

go 1.23.0

replace oracle.com/order => ../order

replace oracle.com/chaos => ../chaos

require (
	oracle.com/chaos v0.0.0-00010101000000-000000000000
	oracle.com/order v0.0.0-00010101000000-000000000000
)
