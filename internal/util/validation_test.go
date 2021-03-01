package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ValidationTestSuite struct {
	suite.Suite
}

func TestInputValidation(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}

func (v *ValidationTestSuite) TestIPValidation() {
	v.T().Run("Test valid IP", func(t *testing.T) {
		assert.True(v.T(), IsIpValid("1.2.3.4"))
	})

	v.T().Run("Test valid IP", func(t *testing.T) {
		assert.True(v.T(), IsIpValid("1.2.3.5"))
	})

	v.T().Run("Test valid IP", func(t *testing.T) {
		assert.True(v.T(), IsIpValid("1.2.3.10"))
	})

	v.T().Run("Test valid IP", func(t *testing.T) {
		assert.True(v.T(), IsIpValid("10.10.9.1"))
	})

	v.T().Run("Test invalid IP 3 digits", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("1.2.3"))
	})

	v.T().Run("Test invalid IP 3 digits with dot", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("1.2.3."))
	})

	v.T().Run("Test invalid IP 3 digits with malformed structure", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("1.2..3"))
	})

	v.T().Run("Test invalid IP 3 digits with Alpha", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("1.2.3.A"))
	})

	v.T().Run("Test invalid IP Alpha with malformed structure", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("A.B..A"))
	})

	v.T().Run("Test invalid IP 4 digits per block", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("1111.2222.3333.4444"))
	})

	v.T().Run("Test invalid IP 4 digits one block", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("1.2222.3.4"))
	})

	v.T().Run("Test valid IP with double quotes", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("\"1.2.3.4\""))
	})

	v.T().Run("Test valid IP with single quotes", func(t *testing.T) {
		assert.False(v.T(), IsIpValid("'1.2.3.4'"))
	})
}

func (v *ValidationTestSuite) TestDomainValidation() {
	v.T().Run("Test valid Domain net", func(t *testing.T) {
		assert.True(v.T(), IsDomainValid("jim.net"))
	})

	v.T().Run("Test valid Domain org", func(t *testing.T) {
		assert.True(v.T(), IsDomainValid("jim.org"))
	})

	v.T().Run("Test valid Domain com", func(t *testing.T) {
		assert.True(v.T(), IsDomainValid("jim.com"))
	})

	v.T().Run("Test invalid domain no TLD", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("jim"))
	})

	v.T().Run("Test invalid domain double dot", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("jim..net"))
	})

	v.T().Run("Test invalid domain with www prefix", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("www.jim.net"))
	})

	v.T().Run("Test invalid domain with http prefix", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("http://jim.net"))
	})

	v.T().Run("Test invalid domain with https prefix", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("http://jim.net"))
	})

	v.T().Run("Test invalid domain with subdomain", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("jim.jims.net"))
	})

	v.T().Run("Test valid domain with double quotes", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("\"jim.net\""))
	})

	v.T().Run("Test valid domain with single quotes", func(t *testing.T) {
		assert.False(v.T(), IsDomainValid("'jim.net'"))
	})
}

func (v *ValidationTestSuite) TestURLValidation() {
	v.T().Run("Test valid URL with http", func(t *testing.T) {
		assert.True(v.T(), IsUrlValid("http://www.jim.com"))
	})

	v.T().Run("Test valid URL with https", func(t *testing.T) {
		assert.True(v.T(), IsUrlValid("https://www.jim.com"))
	})

	v.T().Run("Test valid URL with https and without www", func(t *testing.T) {
		assert.True(v.T(), IsUrlValid("https://jim.com"))
	})

	v.T().Run("Test valid URL with www", func(t *testing.T) {
		assert.True(v.T(), IsUrlValid("www.jim.com"))
	})

	v.T().Run("Test valid URL with http no www", func(t *testing.T) {
		assert.True(v.T(), IsUrlValid("http://nmo.com"))
	})

	v.T().Run("Test invalid URL with httpx", func(t *testing.T) {
		assert.False(v.T(), IsUrlValid("httpx://www.jim.com"))
	})

	v.T().Run("Test invalid URL with htto", func(t *testing.T) {
		assert.False(v.T(), IsUrlValid("htto://www.jim.com"))
	})

	v.T().Run("Test invalid URL with malformed protocol", func(t *testing.T) {
		assert.False(v.T(), IsUrlValid("://www.jim.com"))
	})

	v.T().Run("Test invalid URL with malformed protocol two dots", func(t *testing.T) {
		assert.False(v.T(), IsUrlValid("://www.jim..com"))
	})

	v.T().Run("Test valid URL with two dots", func(t *testing.T) {
		assert.True(v.T(), IsUrlValid("www.jim.co.uk"))
	})

	v.T().Run("Test invalid URL with two dots", func(t *testing.T) {
		assert.False(v.T(), IsUrlValid("www.jim..co.uk"))
	})

	v.T().Run("Test valid URL with double quotes", func(t *testing.T) {
		assert.False(v.T(), IsUrlValid("\"www.jim.co.uk\""))
	})

	v.T().Run("Test valid URL with single quotes", func(t *testing.T) {
		assert.False(v.T(), IsUrlValid("'www.jim.co.uk'"))
	})
}

func (v *ValidationTestSuite) TestRangeValidation() {
	v.T().Run("Test valid Range", func(t *testing.T) {
		assert.True(v.T(), IsRangeValid("1.2.3.4-2.3.4.5"))
	})

	v.T().Run("Test valid Range high range", func(t *testing.T) {
		assert.True(v.T(), IsRangeValid("255.255.255.255-255.200.20.2"))
	})

	v.T().Run("Test valid range", func(t *testing.T) {
		assert.True(v.T(), IsRangeValid("20.2.2.1-20.2.2.3"))
	})

	v.T().Run("Test valid range", func(t *testing.T) {
		assert.True(v.T(), IsRangeValid("192.168.0.0-192.168.255.254"))
	})

	v.T().Run("Test valid range", func(t *testing.T) {
		assert.True(v.T(), IsRangeValid("192.168.0.1-192.168.0.16"))
	})

	v.T().Run("Test valid range", func(t *testing.T) {
		assert.True(v.T(), IsRangeValid("192.168.0.1-192.168.254.254"))
	})

	v.T().Run("Test valid range", func(t *testing.T) {
		assert.True(v.T(), IsRangeValid("192.168.0.1-192.168.16.16"))
	})

	v.T().Run("Test invalid Range too many values", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("1.2.3.4.5-2.3.4.5.9"))
	})

	v.T().Run("Test invalid Range not enough values", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("1.2-3.4"))
	})

	v.T().Run("Test invalid Range not enough values too many dots", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("1.2..-3.4.."))
	})

	v.T().Run("Test invalid Range not enough values too many dots", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("1.2.3.-3.4.5."))
	})

	v.T().Run("Test invalid Range Alphanumeric", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("a.b.c.d-e.f.g.h"))
	})

	v.T().Run("Test invalid Range zeroes", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("0.0.0.0-0.0.0.0"))
	})

	v.T().Run("Test valid Range with double quotes", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("\"1.2.3.4-5.6.7.8\""))
	})

	v.T().Run("Test valid Range with single quotes", func(t *testing.T) {
		assert.False(v.T(), IsRangeValid("'1.2.3.4-5.6.7.8'"))
	})
}

func (v *ValidationTestSuite) TestEmailValidation() {
	v.T().Run("Test valid email firstname", func(t *testing.T) {
		assert.True(v.T(), IsEmailValid("jim@forcepoint.com"))
	})

	v.T().Run("Test valid email first.last", func(t *testing.T) {
		assert.True(v.T(), IsEmailValid("bing.bong@bing.org"))
	})

	v.T().Run("Test valid email double TLD", func(t *testing.T) {
		assert.True(v.T(), IsEmailValid("t@test.co.uk"))
	})

	v.T().Run("Test valid email double TLD first.last", func(t *testing.T) {
		assert.True(v.T(), IsEmailValid("test.case@test.co.uk"))
	})

	v.T().Run("Test valid email - short", func(t *testing.T) {
		assert.True(v.T(), IsEmailValid("t@a.it"))
	})

	v.T().Run("Test valid email - long", func(t *testing.T) {
		assert.True(v.T(), IsEmailValid("purplemonkeydishwasher@thesimpsonswereawfulafterseason10.whatashame"))
	})

	v.T().Run("Test invalid email no prefix", func(t *testing.T) {
		assert.False(v.T(), IsEmailValid("@a.it"))
	})

	v.T().Run("Test invalid email no domain", func(t *testing.T) {
		assert.False(v.T(), IsEmailValid("t@.it"))
	})

	v.T().Run("Test invalid email no domain or dot", func(t *testing.T) {
		assert.False(v.T(), IsEmailValid("t@com"))
	})

	v.T().Run("Test invalid email no domain .com", func(t *testing.T) {
		assert.False(v.T(), IsEmailValid("t@.com"))
	})

	v.T().Run("Test invalid email no domain short", func(t *testing.T) {
		assert.False(v.T(), IsEmailValid("t@.c"))
	})

	v.T().Run("Test invalid email short TLD", func(t *testing.T) {
		assert.False(v.T(), IsEmailValid("t@test.c"))
	})

	v.T().Run("Test invalid email double TLD short", func(t *testing.T) {
		assert.False(v.T(), IsEmailValid("t@test.co.u"))
	})
}
