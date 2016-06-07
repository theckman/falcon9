package f9mission_test

import (
	"testing"

	"github.com/theckman/falcon9/crew"
	"github.com/theckman/falcon9/mission"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	mission *f9mission.Mission
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func addCrew(m *f9mission.Mission, c *C) {
	crew, err := f9crew.NewCrewMember("Jebediah Kerman", "0")
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)

	crew, err = f9crew.NewCrewMember("Bill Kerman", "1")
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)

	crew, err = f9crew.NewCrewMember("Bob Kerman", "2")
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)
}

func (t *TestSuite) SetUpSuite(c *C) {
	mission, err := f9mission.NewMission(
		&f9mission.MissionParams{
			ID:     "42",
			Name:   "Mission: Test Suite",
			GoNoGo: f9mission.VoteQuorum,
		},
	)
	c.Assert(err, IsNil)
	c.Assert(mission, NotNil)

	t.mission = mission
	addCrew(t.mission, c)
}

func (*TestSuite) TestNewMission(c *C) {
	var m *f9mission.Mission
	var err error

	m, err = f9mission.NewMission(nil)
	c.Check(err, ErrorMatches, "mission parameters cannot be nil")
	c.Check(m, IsNil)

	m, err = f9mission.NewMission(&f9mission.MissionParams{})
	c.Check(err, ErrorMatches, "the ID of the mission cannot be an empty string")
	c.Check(m, IsNil)

	m, err = f9mission.NewMission(&f9mission.MissionParams{ID: "id"})
	c.Check(err, IsNil)
	c.Check(m, NotNil)
}

func (t *TestSuite) TestMission_Name(c *C) {
	c.Check(t.mission.Name(), Equals, "Mission: Test Suite")
}

func (t *TestSuite) TestMission_ID(c *C) {
	c.Check(t.mission.ID(), Equals, "42")
}

func (t *TestSuite) TestMission_GNGSetting(c *C) {
	c.Check(t.mission.GNGSetting(), Equals, f9mission.VoteQuorum)
}

func (t *TestSuite) TestMission_Crew(c *C) {
	var crew f9crew.Manifest

	// get the crew
	crew = t.mission.Crew()
	c.Assert(len(crew), Equals, 3)

	// sort the crew
	crew.Sort()

	c.Check(crew[0].Name(), Equals, "Bill Kerman")
	c.Check(crew[0].HashedKey(), Equals, "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")

	c.Check(crew[1].Name(), Equals, "Bob Kerman")
	c.Check(crew[1].HashedKey(), Equals, "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35")

	c.Check(crew[2].Name(), Equals, "Jebediah Kerman")
	c.Check(crew[2].HashedKey(), Equals, "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9")
}

func (t *TestSuite) TestMission_AddCrew(c *C) {
	var err error

	// generate a new crew member
	person, err := f9crew.NewCrewMember("Valentina Kerman", "3")
	c.Assert(err, IsNil)

	// add the crew member
	err = t.mission.AddCrew(person, false)
	c.Assert(err, IsNil)

	// get the crew and sort it
	crew := t.mission.Crew()
	c.Assert(len(crew), Equals, 4)

	crew.Sort()

	c.Check(crew[0].Name(), Equals, "Bill Kerman")
	c.Check(crew[0].HashedKey(), Equals, "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")

	c.Check(crew[1].Name(), Equals, "Bob Kerman")
	c.Check(crew[1].HashedKey(), Equals, "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35")

	c.Check(crew[2].Name(), Equals, "Jebediah Kerman")
	c.Check(crew[2].HashedKey(), Equals, "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9")

	c.Check(crew[3].Name(), Equals, "Valentina Kerman")
	c.Check(crew[3].HashedKey(), Equals, "4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce")

	// assert error if crew present and we aren't replacing them
	err = t.mission.AddCrew(person, false)
	c.Check(err, Equals, f9mission.ErrCrewMemberAlreadyPresent)

	// assert no error if crew present and we are replacing them
	err = t.mission.AddCrew(person, true)
	c.Check(err, IsNil)

	// reset the mission
	t.SetUpSuite(c)
}

func (t *TestSuite) TestMission_RemoveCrew(c *C) {
	var person f9crew.Interface
	var err error

	// remove the person from the crew and make sure it is who we expect
	person, err = t.mission.RemoveCrew("d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35")
	c.Assert(err, IsNil)
	c.Assert(person, NotNil)
	c.Check(person.HashedKey(), Equals, "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35")
	c.Check(person.Name(), Equals, "Bob Kerman")

	// get the crew, and ensure it's the size we expect
	crew := t.mission.Crew()
	c.Check(len(crew), Equals, 2)

	crew.Sort()

	c.Check(crew[0].Name(), Equals, "Bill Kerman")
	c.Check(crew[0].HashedKey(), Equals, "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")

	c.Check(crew[1].Name(), Equals, "Jebediah Kerman")
	c.Check(crew[1].HashedKey(), Equals, "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9")

	// reset the mission
	t.SetUpSuite(c)
}
