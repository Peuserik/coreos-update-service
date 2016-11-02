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

// Test nodes
func (suite storageTestSuite) TestRegisterNode() {
	const machineId = "aaaa-bbbb-ccc-dddd"
	storage := NewLocalDB(path)
	storage.RegisterNode(machineId, "ExTrack", "omd", "123")

	nodes := storage.ListNodes()

	assert.NotNil(suite.T(), nodes[machineId], "Server is listed")
}

// Test a version version update
func (suite storageTestSuite) TestUpdateVersion() {
	const versionId = "1.2.3"
	storage := NewLocalDB(path)
	version := CoreOSVersion{
		VersionId: versionId,
		URL: "http://test_url.com",
		Name: "update.gz",
		Hash: "abcdef",
		Signature: "abcdef",
		Size: 2000,
	}
	storage.UpdateVersion(version)

	version2 := storage.GetVersion(versionId)
	assert.NotNil(suite.T(), version2)
	assert.Equal(suite.T(), version, version2, "Storage version and reload aren't equal")
}

// Test update tracks
func (suite storageTestSuite) TestTrackUpdate() {
	storage := NewLocalDB(path)
	tracks := map[string]string {
		"TestTrack": "1.2.3",
	}
	storage.UpdateTracks(tracks)

	tracks2 := storage.GetTracks()
	assert.NotNil(suite.T(), tracks2)
	assert.Equal(suite.T(), tracks, tracks2, "Tracks storage and reload aren't equal")
}


func (suite storageTestSuite) TestLogEvent() {
	storage := NewLocalDB(path)
	storage.LogEvent("aaaa-bbbb-ccc-dddd", "", "", "")
	assert.True(suite.T(), true)
}


// Remove level db directory after each test execution
func (suite storageTestSuite) TearDownTest() ()  {
	os.RemoveAll(path)
}

func TestStorageTestSuite(t *testing.T) {
	tests := new(storageTestSuite)
	suite.Run(t, tests)
}
