package testutils

import (
	"gorm.io/gorm"
	"testing"
)

// MockGormDB creates a mock for the global db variable
// This returns a mock implementation that can be used in place of a real *gorm.DB
func MockGormDB(t *testing.T) *gorm.DB {
	// Create a wrapper that forwards calls to the mock
	db := &gorm.DB{
		Config: &gorm.Config{},
	}
	
	return db
}