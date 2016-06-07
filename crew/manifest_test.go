package f9crew_test

import (
	"github.com/theckman/falcon9/crew"
	. "gopkg.in/check.v1"
)

func (*TestSuite) TestManifest_Len(c *C) {
	m := make(f9crew.Manifest, 2, 3)

	crew, err := f9crew.NewCrewMember("TestUser1", "42")
	c.Assert(err, IsNil)

	m[0] = crew

	crew, err = f9crew.NewCrewMember("TestUser2", "33")
	c.Assert(err, IsNil)

	m[1] = crew

	c.Check(m.Len(), Equals, 2)

	crew, err = f9crew.NewCrewMember("TestUser3", "key")
	c.Assert(err, IsNil)

	m = append(m, crew)

	c.Check(m.Len(), Equals, 3)
}

func (*TestSuite) TestManifest_Swap(c *C) {
	m := make(f9crew.Manifest, 3)

	crew, err := f9crew.NewCrewMember("TestUser1", "42")
	c.Assert(err, IsNil)

	m[0] = crew

	crew, err = f9crew.NewCrewMember("TestUser2", "33")
	c.Assert(err, IsNil)

	m[1] = crew

	crew, err = f9crew.NewCrewMember("TestUser3", "key")
	c.Assert(err, IsNil)

	m[2] = crew

	old0, old2 := m[0], m[2]

	m.Swap(0, 2)
	c.Check(m[0], DeepEquals, old2)
	c.Check(m[2], DeepEquals, old0)
}

func (*TestSuite) TestManifest_Less(c *C) {
	m := make(f9crew.Manifest, 3)

	crew, err := f9crew.NewCrewMember("TestUser1", "42")
	c.Assert(err, IsNil)

	m[0] = crew

	crew, err = f9crew.NewCrewMember("TestUser2", "vdc")
	c.Assert(err, IsNil)

	m[1] = crew

	crew, err = f9crew.NewCrewMember("TestUser2", "abc")
	c.Assert(err, IsNil)

	m[2] = crew

	c.Check(m.Less(0, 1), Equals, true)
	c.Check(m.Less(1, 2), Equals, false) // uses HashedKey to sort
}

func (*TestSuite) TestManifest_Sort(c *C) {
	m := make(f9crew.Manifest, 3)

	crew, err := f9crew.NewCrewMember("TestUser4", "42")
	c.Assert(err, IsNil)

	m[0] = crew

	crew, err = f9crew.NewCrewMember("TestUser2", "vdc")
	c.Assert(err, IsNil)

	m[1] = crew

	crew, err = f9crew.NewCrewMember("TestUser2", "abc")
	c.Assert(err, IsNil)

	m[2] = crew

	m.Sort()

	c.Check(m[0].Name(), Equals, "TestUser2")
	c.Check(m[0].HashedKey(), Equals, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad")

	c.Check(m[1].Name(), Equals, "TestUser2")
	c.Check(m[1].HashedKey(), Equals, "e54a5df6ff5bc3208e2c9ebd47471ae84e3ce0a64d912d6a81af7e4f477d5df8")

	c.Check(m[2].Name(), Equals, "TestUser4")
	c.Check(m[2].HashedKey(), Equals, "73475cb40a568e8da8a045ced110137e159f890ac4da883b6b17dc651b3a8049")
}
