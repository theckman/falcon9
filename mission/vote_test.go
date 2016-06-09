package f9mission_test

import (
	"github.com/theckman/falcon9/mission"
	. "gopkg.in/check.v1"
)

func (*TestSuite) TestVote_String(c *C) {
	c.Check(f9mission.VoteAbstain.String(), Equals, "Abstain")
	c.Check(f9mission.VoteNo.String(), Equals, "No")
	c.Check(f9mission.VoteYes.String(), Equals, "Yes")
	c.Check(f9mission.Vote(100).String(), Equals, "Unknown")
}
