// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package f9crew_test

import (
	"testing"

	"github.com/theckman/falcon9/crew"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	crew *f9crew.CrewMember
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (t *TestSuite) SetUpSuite(c *C) {
	cm, err := f9crew.NewCrewMember("Test Case", "anothertest")
	c.Assert(err, IsNil)
	c.Assert(cm, NotNil)

	t.crew = cm
}

func (*TestSuite) TestHashKey(c *C) {
	tests := []struct {
		i, o string
	}{
		{"test", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
		{"anothertest", "9bb5bde1a740465d012231e350aa8934f64d078ac1349d6a10852cbf1369d15f"},
	}

	for _, test := range tests {
		c.Check(f9crew.HashKey(test.i), Equals, test.o)
	}
}

func (*TestSuite) BenchmarkHashKey(c *C) {
	for i := 0; i < c.N; i++ {
		f9crew.HashKey("testing a unique value")
	}
}

func (*TestSuite) TestNewCrewMember(c *C) {
	var cm *f9crew.CrewMember
	var err error

	cm, err = f9crew.NewCrewMember("", "key")
	c.Assert(err, NotNil)
	c.Check(cm, IsNil)
	c.Check(err.Error(), Equals, "the crew member's name cannot be an empty value")

	cm, err = f9crew.NewCrewMember("name", "")
	c.Assert(err, NotNil)
	c.Check(cm, IsNil)
	c.Check(err.Error(), Equals, "the crew member's key cannot be an empty value")

	cm, err = f9crew.NewCrewMember("name", "key")
	c.Assert(err, IsNil)
	c.Assert(cm, NotNil)
}

func (t *TestSuite) TestCrewMember_Name(c *C) {
	c.Check(t.crew.Name(), Equals, "Test Case")
}

func (t *TestSuite) TestCrewMember_HashedKey(c *C) {
	c.Check(t.crew.HashedKey(), Equals, "9bb5bde1a740465d012231e350aa8934f64d078ac1349d6a10852cbf1369d15f")
}
