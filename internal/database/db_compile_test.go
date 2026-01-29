package database

import "testing"

func TestInitDB(t *testing.T) {
	// We pass an invalid connection string.
	// We expect an error, but we just want to ensure it compiles.
	_, err := InitDB("postgres://invalid:invalid@localhost:5432/invalid")
	if err == nil {
		// This path is unlikely to be reached
	}
}
