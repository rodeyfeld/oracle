package order

import "testing"

func TestConnectPostgres(t *testing.T) {
	conn := ConnectPostgres()
	print(conn)
}
