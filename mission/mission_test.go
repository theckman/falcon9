package f9mission_test

import (
	"testing"
	"time"

	"github.com/theckman/falcon9/crew"
	"github.com/theckman/falcon9/mission"
	"github.com/theckman/go-fsm"

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
			GoNoGo: f9mission.GNGQuorum,
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
	c.Check(t.mission.GNGSetting(), Equals, f9mission.GNGQuorum)
}

func (t *TestSuite) TestMission_CurrentState(c *C) {
	// reset the mission
	defer t.SetUpSuite(c)

	//
	// Test that the initial state is StateReady
	//
	c.Check(t.mission.CurrentState(), Equals, f9mission.StateReady)

	//
	// Test that the state changing is seen
	//
	c.Assert(t.mission.Initiate(), Equals, nil)
	c.Check(t.mission.CurrentState(), Equals, f9mission.StateVoting)
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

	// reset the mission
	defer t.SetUpSuite(c)

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

	//
	// Test that AddCrew will abort a blastoff
	//
	err = t.mission.Initiate()
	c.Assert(err, IsNil)

	ok, err := t.mission.UpdateVote("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, false)

	ok, err = t.mission.UpdateVote("5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, false)

	ok, err = t.mission.UpdateVote("4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
}

func (t *TestSuite) TestMission_RemoveCrew(c *C) {
	var person f9crew.Interface
	var err error

	// reset the mission
	defer t.SetUpSuite(c)

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
}

func (*TestSuite) TestMission_Initiate(c *C) {
	var m *f9mission.Mission
	var err error
	var state fsm.State

	m, err = f9mission.NewMission(&f9mission.MissionParams{
		ID:     "id",
		GoNoGo: f9mission.GNGQuorum,
	})
	c.Check(err, IsNil)
	c.Check(m, NotNil)

	err = m.Initiate()
	c.Check(err, Equals, f9mission.ErrNoAssignedCrew)

	crew, err := f9crew.NewCrewMember("Jebediah Kerman", "0")
	c.Assert(err, IsNil)

	err = m.AddCrew(crew, false)
	c.Assert(err, IsNil)

	err = m.Initiate()
	c.Check(err, IsNil)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateVoting)

	err = m.Initiate()
	c.Check(err, Equals, f9mission.ErrMissionInProgress)
}

func (t *TestSuite) TestMission_UpdateVote(c *C) {
	var bol bool
	var err error
	var state fsm.State

	m, err := f9mission.NewMission(&f9mission.MissionParams{
		ID:                  "id",
		GoNoGo:              f9mission.GNGQuorum,
		BlastoffingCooldown: time.Millisecond * 100,
	})
	c.Check(err, IsNil)
	c.Check(m, NotNil)

	crew, err := f9crew.NewCrewMember("Jebediah Kerman", "0")
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)

	crew, err = f9crew.NewCrewMember("Bill Kerman", "1")
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)

	err = m.Initiate()
	c.Assert(err, IsNil)

	//
	// Test that, with GNGQuorum, less than three crew members falls
	// back to GNGAll mode
	//
	bol, err = m.UpdateVote("5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(bol, Equals, false)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateVoting)

	bol, err = m.UpdateVote("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(bol, Equals, true)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateBlastoffing)

	//
	// Test that adding third crew member, with GNGQuorum and two votes,
	// marks the vote as being successful and aborts
	//
	crew, err = f9crew.NewCrewMember("Bob Kerman", "2")
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)

	_, bol = m.Tally()
	c.Check(bol, Equals, true)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateAborted)

	//
	// Test that adding more crew, thus changing quorum, causes the
	// vote to fall back to a fail and aborts the mission.
	//
	crew, err = f9crew.NewCrewMember("Valentina Kerman", "3")
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)

	crew, err = f9crew.NewCrewMember("Dildo Kerman", "4") /* lol, I hope someone gets this reference */
	c.Assert(err, IsNil)
	c.Assert(m.AddCrew(crew, false), IsNil)

	_, bol = m.Tally()
	c.Check(bol, Equals, false)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateAborted)

	//
	// RESTART VOTING AFTER PLANNED ABORT
	//
	err = m.Initiate()
	c.Assert(err, IsNil)

	// add votes back
	bol, err = m.UpdateVote("5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(bol, Equals, false)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateVoting)

	bol, err = m.UpdateVote("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(bol, Equals, false)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateVoting)

	//
	// Test that updating a cast vote doesn't change anything
	//
	bol, err = m.UpdateVote("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(bol, Equals, false)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateVoting)

	//
	// Test that adding another Yes vote gets us quorum.
	//
	bol, err = m.UpdateVote("4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce", f9mission.VoteYes)
	c.Assert(err, IsNil)
	c.Check(bol, Equals, true)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateBlastoffing)

	//
	// Test that adding an abort vote aborts
	//
	bol, err = m.UpdateVote("d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35", f9mission.VoteAbort)
	c.Assert(err, IsNil)
	c.Check(bol, Equals, false)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateAborted)

	//
	// RESTART VOTING AFTER PLANNED ABORT
	//
	err = m.Initiate()
	c.Assert(err, IsNil)

	// add votes back
	_, err = m.UpdateVote("5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9", f9mission.VoteYes)
	c.Assert(err, IsNil)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateVoting)

	_, err = m.UpdateVote("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", f9mission.VoteYes)
	c.Assert(err, IsNil)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateVoting)

	_, err = m.UpdateVote("4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce", f9mission.VoteYes)
	c.Assert(err, IsNil)

	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateBlastoffing)

	//
	// Test that state transitions to StateFinished after the BlastoffingCooldown timer
	//

	time.Sleep(time.Millisecond * 110)
	state = m.CurrentState()
	c.Check(state, Equals, f9mission.StateFinished)
}

func (t *TestSuite) TestMission_Tally(c *C) {
	var rdy bool
	var err error
	var tally f9mission.Tally

	// reset the mission
	defer t.SetUpSuite(c)

	//
	// Test that nil is returned before Initiate() function is called
	//
	tally, rdy = t.mission.Tally()
	c.Assert(len(tally), Equals, 0)
	c.Check(tally, DeepEquals, f9mission.Tally(nil))
	c.Check(rdy, Equals, false)

	err = t.mission.Initiate()
	c.Assert(err, Equals, nil)

	//
	// Test that not nil is returned after Initiate() function is called
	//
	tally, rdy = t.mission.Tally()
	c.Assert(len(tally), Equals, 0)
	c.Check(tally, Not(DeepEquals), f9mission.Tally(nil))
	c.Check(rdy, Equals, false)

	//
	// Test that Tally() can do maths
	//
	_, err = t.mission.UpdateVote("5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9", f9mission.VoteYes)
	c.Assert(err, IsNil)

	tally, rdy = t.mission.Tally()
	c.Assert(len(tally), Equals, 1)
	c.Check(rdy, Equals, false)
	c.Check(tally[f9mission.VoteYes], Equals, 1)
	c.Check(tally[f9mission.VoteNo], Equals, 0)
	c.Check(tally[f9mission.VoteAbstain], Equals, 0)
	c.Check(tally[f9mission.VoteAbort], Equals, 0)

	_, err = t.mission.UpdateVote("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", f9mission.VoteNo)
	c.Assert(err, IsNil)

	tally, rdy = t.mission.Tally()
	c.Assert(len(tally), Equals, 2)
	c.Check(rdy, Equals, false)
	c.Check(tally[f9mission.VoteYes], Equals, 1)
	c.Check(tally[f9mission.VoteNo], Equals, 1)
	c.Check(tally[f9mission.VoteAbstain], Equals, 0)
	c.Check(tally[f9mission.VoteAbort], Equals, 0)

	_, err = t.mission.UpdateVote("d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35", f9mission.VoteAbstain)
	c.Assert(err, IsNil)

	tally, rdy = t.mission.Tally()
	c.Assert(len(tally), Equals, 3)
	c.Check(rdy, Equals, false)
	c.Check(tally[f9mission.VoteYes], Equals, 1)
	c.Check(tally[f9mission.VoteNo], Equals, 1)
	c.Check(tally[f9mission.VoteAbstain], Equals, 1)
	c.Check(tally[f9mission.VoteAbort], Equals, 0)

	//
	// Test that Tally() marks rdy as true when ready
	//
	_, err = t.mission.UpdateVote("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", f9mission.VoteYes)
	c.Assert(err, IsNil)

	tally, rdy = t.mission.Tally()
	c.Assert(len(tally), Equals, 2)
	c.Check(rdy, Equals, true)
	c.Check(tally[f9mission.VoteYes], Equals, 2)
	c.Check(tally[f9mission.VoteNo], Equals, 0)
	c.Check(tally[f9mission.VoteAbstain], Equals, 1)
	c.Check(tally[f9mission.VoteAbort], Equals, 0)

	//
	// Test that Tally() marks rdy as false when there is an abort
	//
	_, err = t.mission.UpdateVote("d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35", f9mission.VoteAbort)
	c.Assert(err, IsNil)

	tally, rdy = t.mission.Tally()
	c.Assert(len(tally), Equals, 2)
	c.Check(rdy, Equals, false)
	c.Check(tally[f9mission.VoteYes], Equals, 2)
	c.Check(tally[f9mission.VoteNo], Equals, 0)
	c.Check(tally[f9mission.VoteAbstain], Equals, 0)
	c.Check(tally[f9mission.VoteAbort], Equals, 1)
}
