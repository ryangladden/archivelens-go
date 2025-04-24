package db

import (
	"testing"
)

func TestConnect(t *testing.T) {
	err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
}
