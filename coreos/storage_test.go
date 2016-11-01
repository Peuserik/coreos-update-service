package coreos

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const path = "._local_db_test"

type storageTestSuite struct {
	suite.Suite
}

// Test that the initialization work and it creates a level db directory
func (suite storageTestSuite) TestLocalDBCreatesFile() {
	NewLocalDB(path)
	info , err := os.Stat(path)
	assert.False(suite.T(), os.IsNotExist(err), "Local storage creates a directory")
	assert.True(suite.T(), info.IsDir(), "Local storage creates a directory")
}

// Remove level db directory after each test execution
func (suite storageTestSuite) TearDownTest() ()  {
	os.RemoveAll(path)
}

func TestStorageTestSuite(t *testing.T) {
	tests := new(storageTestSuite)
	suite.Run(t, tests)
}