package storage_test

import (
	"testing"
	"time"

	"github.com/fabinzne/go-simple-redis/internal/storage"
	"github.com/stretchr/testify/suite"
)

type StorageSetTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageSetTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageSetTestSuite) TestSet() {
	s.dataStore.Set("key", "value")
	value, ok := s.dataStore.Get("key")
	s.True(ok)
	s.Equal("value", value)
}

type StorageGetTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageGetTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageGetTestSuite) TestGet() {
	s.dataStore.Set("key", "value")
	value, ok := s.dataStore.Get("key")
	s.True(ok)
	s.Equal("value", value)
}

func (s *StorageGetTestSuite) TestGetNotFound() {
	value, ok := s.dataStore.Get("key")
	s.False(ok)
	s.Nil(value)
}

type StorageDeleteTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageDeleteTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageDeleteTestSuite) TestDelete() {
	s.dataStore.Set("key", "value")
	s.dataStore.Delete("key")
	value, ok := s.dataStore.Get("key")
	s.False(ok)
	s.Nil(value)
}

type StorageSetWithTTLTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageSetWithTTLTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageSetWithTTLTestSuite) TestSetWithTTL() {
	s.dataStore.SetWithTTL("key", "value", 1)
	value, ok := s.dataStore.Get("key")
	s.True(ok)
	s.Equal("value", value)
}

func (s *StorageSetWithTTLTestSuite) TestSetWithTTLExpired() {
	s.dataStore.SetWithTTL("key", "value", 1)
	time.Sleep(2 * time.Second)
	value, ok := s.dataStore.Get("key")
	s.False(ok)
	s.Nil(value)
}

type StorageExpireTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageExpireTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageExpireTestSuite) TestExpire() {
	s.dataStore.Set("key", "value")
	s.True(s.dataStore.Expire("key", 1))
}

func (s *StorageExpireTestSuite) TestExpireNotFound() {
	s.False(s.dataStore.Expire("key", 1))
}

func (s *StorageExpireTestSuite) TestExpireExpired() {
	s.dataStore.Set("key", "value")
	time.Sleep(2 * time.Second)
	s.True(s.dataStore.Expire("key", 1))
}

func (s *StorageExpireTestSuite) TestExpireWithTTLExpired() {
	s.dataStore.SetWithTTL("key", "value", 1)
	time.Sleep(2 * time.Second)
	s.False(s.dataStore.Expire("key",1))
}

func (s *StorageExpireTestSuite) TestExpireWithTTLNotFound() {
	s.False(s.dataStore.Expire("key", 1))
}

func (s *StorageExpireTestSuite) TestExpireWithTTLNotExpired() {
	s.dataStore.SetWithTTL("key", "value", 2)
	s.True(s.dataStore.Expire("key", 1))
}

type StorageFlushTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageFlushTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageFlushTestSuite) TestFlush() {
	s.dataStore.Set("key", "value")
	s.dataStore.Flush()
	value, ok := s.dataStore.Get("key")
	s.False(ok)
	s.Nil(value)
}

func (s *StorageFlushTestSuite) TestFlushWithTTL() {
	s.dataStore.SetWithTTL("key", "value", 1)
	s.dataStore.Flush()
	value, ok := s.dataStore.Get("key")
	s.False(ok)
	s.Nil(value)
}

type StorageSafeGetTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageSafeGetTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageSafeGetTestSuite) TestSafeGet() {
	s.dataStore.Set("key", "value")
	value, ok, _ := s.dataStore.SafeGetKeyInfo("key")
	s.True(ok)
	s.Equal("value", value)
}

func (s *StorageSafeGetTestSuite) TestSafeGetWithTTL() {
	s.dataStore.SetWithTTL("key", "value", 1)
	value, ok, exp := s.dataStore.SafeGetKeyInfo("key")
	s.True(ok)
	s.Equal("value", value)
	s.NotNil(exp)
}

func (s *StorageSafeGetTestSuite) TestSafeGetNotFound() {
	value, ok, exp := s.dataStore.SafeGetKeyInfo("key")
	s.False(ok)
	s.Nil(value)
	s.Equal(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), exp)
}

func (s *StorageSafeGetTestSuite) TestSafeGetExpired() {
	s.dataStore.SetWithTTL("key", "value", 1)
	time.Sleep(2 * time.Second)
	value, ok, _ := s.dataStore.SafeGetKeyInfo("key")
	s.False(ok)
	s.Nil(value)
}

type StorageSaveToFileTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageSaveToFileTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageSaveToFileTestSuite) TestSaveToFile() {
	s.dataStore.Set("key", "value")
	err := s.dataStore.SaveToFile("test.dump")
	s.Nil(err)
}

func (s *StorageSaveToFileTestSuite) TestSaveToFileWithTTL() {
	s.dataStore.SetWithTTL("key", "value", 1)
	err := s.dataStore.SaveToFile("test.dump")
	s.Nil(err)
}

type StorageLoadFromFileTestSuite struct {
	suite.Suite
	dataStore *storage.DataStore
}

func (s *StorageLoadFromFileTestSuite) SetupTest() {
	s.dataStore = storage.NewDataStore()
}

func (s *StorageLoadFromFileTestSuite) TestLoadFromFile() {
	s.dataStore.Set("key", "value")
	err := s.dataStore.SaveToFile("test.dump")
	s.Nil(err)

	s.dataStore.Flush()
	err = s.dataStore.LoadFromFile("test.dump")
	s.Nil(err)

	value, ok := s.dataStore.Get("key")
	s.True(ok)
	s.Equal("value", value)
}

func (s *StorageLoadFromFileTestSuite) TestLoadFromFileWithTTL() {
	s.dataStore.SetWithTTL("key", "value", 1)
	err := s.dataStore.SaveToFile("test.dump")
	s.Nil(err)

	s.dataStore.Flush()
	err = s.dataStore.LoadFromFile("test.dump")
	s.Nil(err)

	value, ok := s.dataStore.Get("key")
	s.True(ok)
	s.Equal("value", value)
}


func TestStorage(t *testing.T){
	suite.Run(t, new(StorageSetTestSuite))
	suite.Run(t, new(StorageGetTestSuite))
	suite.Run(t, new(StorageDeleteTestSuite))
	suite.Run(t, new(StorageSetWithTTLTestSuite))
	suite.Run(t, new(StorageExpireTestSuite))
	suite.Run(t, new(StorageFlushTestSuite))
	suite.Run(t, new(StorageSafeGetTestSuite))
	suite.Run(t, new(StorageSaveToFileTestSuite))
	suite.Run(t, new(StorageLoadFromFileTestSuite))
}