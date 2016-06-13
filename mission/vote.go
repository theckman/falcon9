package f9mission

// Vote is the type for someone's vote
type Vote uint8

const (
	// VoteAbstain is the default Vote value. It's for when a person has not
	// voted yet. Abstaining is treated the same as a no. However, this is here
	// to try and avoid breaking API compatibility in the future.
	VoteAbstain Vote = iota

	// VoteNo is the no Vote. It means the client is not ready.
	VoteNo

	// VoteYes is the yes Vote. The crew member is ready for blastoff!
	VoteYes

	// VoteAbort is the vote to abort. This is only available for use after the
	// countdown has began. Depending on your mission parameters, a single abort
	// may scrub the launch.
	VoteAbort
)

func (v Vote) String() string {
	switch v {
	case VoteAbstain:
		return "Abstain"
	case VoteNo:
		return "No"
	case VoteYes:
		return "Yes"
	case VoteAbort:
		return "Abort"
	default:
		return "Unknown"
	}
}
