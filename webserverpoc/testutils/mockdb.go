package testutils

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB provides a common mock implementation for database operations in tests
type MockDB struct {
	mock.Mock
}

// First mocks the GORM First method
func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

// Create mocks the GORM Create method
func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// Save mocks the GORM Save method
func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// Find mocks the GORM Find method
func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

// Where mocks the GORM Where method
func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

// Preload mocks the GORM Preload method
func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}

// Delete mocks the GORM Delete method
func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(value, conds)
	return args.Get(0).(*gorm.DB)
}

// Select mocks the GORM Select method
func (m *MockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

// Joins mocks the GORM Joins method
func (m *MockDB) Joins(query string, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

// Begin mocks the GORM Begin method to start a transaction
func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// Commit mocks the GORM Commit method to commit a transaction
func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// Rollback mocks the GORM Rollback method to rollback a transaction
func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// Model mocks the GORM Model method
func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// Count mocks the GORM Count method
func (m *MockDB) Count(count *int64) *gorm.DB {
	args := m.Called(count)
	return args.Get(0).(*gorm.DB)
}

// CreateMockDB creates a new MockDB instance with a success result for chaining
func CreateMockDB() (*MockDB, *gorm.DB) {
	mockDB := new(MockDB)
	mockResult := &gorm.DB{Error: nil}
	return mockDB, mockResult
}

// CreateMockDBWithError creates a new MockDB instance with an error result for chaining
func CreateMockDBWithError(err error) (*MockDB, *gorm.DB) {
	mockDB := new(MockDB)
	mockResult := &gorm.DB{Error: err}
	return mockDB, mockResult
}

// NewMockGorm creates a mocked GORM DB instance
func NewMockGorm() *gorm.DB {
	// This is just a placeholder to return a *gorm.DB without actual database connection
	// In real tests, you'd use this with sqlmock or other mocking approach
	return &gorm.DB{}
}