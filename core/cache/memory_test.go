package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MemoryTestSuite struct {
	suite.Suite
	Mc *MemoryCache
}

func (s *MemoryTestSuite) SetupTest() {
	s.Mc = Boot()
}

func (s *MemoryTestSuite) TearDownAllSuite() {
	s.Mc = nil
}

func TestRunMemoryTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}

func (s *MemoryTestSuite) TestCacheBoot() {
	assert.NotNil(s.T(), s.Mc)
}

func (s *MemoryTestSuite) TestAddToCache() {
	s.Mc.Add("foo", "bar")

	assert.NotNil(s.T(), s.Mc)
}

func (s *MemoryTestSuite) TestGetToCache() {
	s.Mc.Add("foo", "bar")

	data := s.Mc.Get("foo")
	assert.Equal(s.T(), "bar", data)
}

func (s *MemoryTestSuite) TestDeleteByInCache() {
	s.Mc.Add("foo", "bar")

	data := s.Mc.Get("foo")
	assert.Equal(s.T(), "bar", data)

	isDebited := s.Mc.DeleteBy("foo")
	assert.True(s.T(), isDebited)
}

func (s *MemoryTestSuite) TestGetKeyNotInCache() {
	data := s.Mc.Get("some")

	assert.Nil(s.T(), data)
}
