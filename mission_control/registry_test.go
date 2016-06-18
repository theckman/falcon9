package f9missioncontrol_test

import (
	"fmt"
	"math/rand"

	"github.com/theckman/falcon9/mission"
	"github.com/theckman/falcon9/mission_control"

	. "gopkg.in/check.v1"
)

const maxUint16 = int(^uint16(0))

func randUint16() uint16 {
	return uint16(rand.Intn(maxUint16))
}

func tearDownRegistry(c *C) {
	ids := f9missioncontrol.ListMissions()

	for _, id := range ids {
		mission := f9missioncontrol.RemoveMission(id)
		c.Check(mission, NotNil)
	}
}

func (*TestSuite) TestAddMission(c *C) {
	var err error

	// clean up the registry
	defer tearDownRegistry(c)

	mission := &f9mission.Mission{}
	id := randUint16()

	err = f9missioncontrol.AddMission(id, mission)
	c.Assert(err, IsNil)

	mIfc := f9missioncontrol.GetMission(id)
	c.Check(mIfc, NotNil)

	//
	// Test that you can't register it twice
	//
	err = f9missioncontrol.AddMission(id, mission)
	c.Assert(err, NotNil)
}

func (*TestSuite) TestListMissions(c *C) {
	var missions []uint16

	// clean up the registry
	defer tearDownRegistry(c)

	set := make(map[uint16]struct{})

	// set up the registry for this test
	for i := 0; i < 5; i++ {
		id := randUint16()

		mp := &f9mission.MissionParams{
			ID:   fmt.Sprintf("testID-%d", id),
			Name: fmt.Sprintf("testName-%d", id),
		}

		mission, err := f9mission.NewMission(mp)
		c.Assert(err, IsNil)

		err = f9missioncontrol.AddMission(id, mission)
		c.Assert(err, IsNil)

		set[id] = struct{}{}
	}

	missions = f9missioncontrol.ListMissions()
	c.Assert(len(missions), Equals, 5)

	for _, id := range missions {
		_, ok := set[id]
		c.Assert(ok, Equals, true)

		missionIfc := f9missioncontrol.GetMission(id)
		c.Check(missionIfc.Name(), Equals, fmt.Sprintf("testName-%d", id))
		c.Check(missionIfc.ID(), Equals, fmt.Sprintf("testID-%d", id))
	}
}

func (*TestSuite) TestGetMission(c *C) {
	var mIfc f9mission.Interface
	var err error

	// clean up the registry
	defer tearDownRegistry(c)

	mp := &f9mission.MissionParams{
		ID:   "testID",
		Name: "testName",
	}

	mission, err := f9mission.NewMission(mp)
	c.Assert(err, IsNil)

	id := randUint16()

	err = f9missioncontrol.AddMission(id, mission)
	c.Assert(err, IsNil)

	//
	// Test when mission exists
	//
	mIfc = f9missioncontrol.GetMission(id)
	c.Assert(mIfc, NotNil)
	c.Check(mIfc.Name(), Equals, "testName")
	c.Check(mIfc.ID(), Equals, "testID")

	//
	// Test when mission does not exist
	//
	mIfc = f9missioncontrol.GetMission(id + 1)
	c.Assert(mIfc, IsNil)
}

func (*TestSuite) TestRemoveMission(c *C) {
	var err error

	// clean up the registry
	defer tearDownRegistry(c)

	mp := &f9mission.MissionParams{
		ID:   "testID",
		Name: "testName",
	}

	mission, err := f9mission.NewMission(mp)
	c.Assert(err, IsNil)

	id := randUint16()

	err = f9missioncontrol.AddMission(id, mission)
	c.Assert(err, IsNil)

	c.Check(len(f9missioncontrol.ListMissions()), Equals, 1)

	//
	// Test when mission does not exist
	//
	mIfc := f9missioncontrol.RemoveMission(id + 1)
	c.Check(mIfc, IsNil)

	//
	// Test when mission exists
	//
	mIfc = f9missioncontrol.RemoveMission(id)
	c.Assert(mIfc, NotNil)
	c.Check(mIfc.Name(), Equals, "testName")
	c.Check(mIfc.ID(), Equals, "testID")
	c.Check(len(f9missioncontrol.ListMissions()), Equals, 0)
}
