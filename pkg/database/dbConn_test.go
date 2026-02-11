package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	// Set the database URL for the test
	os.Setenv("DB_URL", "postgres://postgres:postgres@localhost:5432/logistic_db?sslmode=disable")

	// Attempt to connect to the database
	db, err := ConnectDB().DB()
	assert.NoError(t, err, "Failed to connect to database")
	assert.NotNil(t, db, "Database connection should not be nil")
}

