package f9mission

import (
	"errors"
	"sync"

	"github.com/theckman/falcon9/crew"
)

// GNGSetting is the type that defines the behavior of the mission parameters for
// blastoff. When the Go/No-Go call is made, how many "Go"s are required for
// blastoff. The default value is to require all crew members vote Go.
type GNGSetting uint8

const (
	// GNGAll is the GoNoGo setting for requiring that all crew members vote Go.
	GNGAll GNGSetting = iota

	// GNGQuorum is the GoNoGo setting for requiring that only a quorum number
	// of crew members vote for blastoff. If there are too few crew members to
	// reach quorum without all voting "Go", it falls back to GNGAll mode.
	GNGQuorum
)

// Tally is the map used for the current voting result tally. The key is the Vote kind.
type Tally map[Vote]int

// Results is the type for voting results per the crew member's answer. The key is
// the crew member's HashedKey.
type Results map[string]Vote

var errUseNewMission = errors.New("use f9mission.NewMission() to create the *Mission struct")

// ErrCrewMemberAlreadyPresent is the error returned from *Mission.AddUser
// if the crew member is already present within the mission. Ths consumer
// may not consider this an error condition.
var ErrCrewMemberAlreadyPresent = errors.New("the user you are trying to add already exists")

// ErrCrewMemberNotPresent is the error returned from RemoveCrew() if the crew member
// is not assigned to this mission.
var ErrCrewMemberNotPresent = errors.New("the crew member you've tried to remove is not present")

// ErrNoAssignedCrew is the error returned from Initiate() if a Go/No-Go was started
// without there being any crew members assigned to the mission.
var ErrNoAssignedCrew = errors.New("before Initiating a Go/No-Go the mission must have crew assigned")

// ErrGNGInProgress is the error returned from Initiate() if a Go/No-Go is in progress.
var ErrGNGInProgress = errors.New("Go/No-Go vote is currently in progress")

// ErrGNGNotInProgress is the error returned from UpdateVote() if a Go/No-Go is not in progress
var ErrGNGNotInProgress = errors.New("Go/No-Go vote is *NOT* currently in progress")

// InterfaceManageCrew is the mission-specific interface for adding crew members.
type InterfaceManageCrew interface {
	// AddCrew is a function to add a new crew member to this mission. If the crew
	// member already exists (identified by their HashedKey), and replace is set to
	// false, this will return an f9crew.ErrCrewMemberAlreadyPresent error. However,
	// if replace is true it will simply replace the crew member with the new one.
	//
	// The latter is useful for a client quickly rejoining the session
	// after a network interruption. This assume clients have a unique key.
	AddCrew(crew f9crew.Interface, replace bool) error

	// RemoveCrew is function to remove a crew member from the mission.
	// The crew member's HashedKey is used to do the lookup for determining which
	// crew member to remove from the mission. This returns the crew member being
	// removed, if a consumer wishes to use it.
	RemoveCrew(hashedKey string) (f9crew.Interface, error)

	// Crew is a function that returns an f9crew.Manifest. This is a
	// representation of the crew for the current mission. This function
	// returns the values unsorted, but the returns value will have a
	// Sort() method.
	Crew() f9crew.Manifest
}

// InterfaceLaunchControl is the interface to the launch control systems.
// This includes Go/No-Go voting
type InterfaceLaunchControl interface {
	// Initiate is the function used to start the launch sequence.
	// By calling this function the mission enters the go/no-go state.
	// The five-second countdown will begin once the proper number of
	// "go" votes have been reached.
	Initiate() error

	// UpdateVote updates the vote of a crew member for the current mission.
	// The bool value returned indicates whether there have been enough "Go"
	// votes to proceed with blastoff.
	//
	// If the mission is not initialized this will return a ErrGNGNotInProgress
	// error. If the crew member is not assigned to this mission, this will return
	// a ErrCrewMembeverNotPresent error.
	UpdateVote(hashedKey string, vote Vote) (bool, error)

	// Tally returns the tally of votes and whether there are enough votes
	// to proceed with the mission.
	Tally() (Tally, bool)
}

// InterfaceAccessors is an interface type for accessor methods of the
// mission parameters.
type InterfaceAccessors interface {
	// Name returns the name of the mission. This is an optional value, so it
	// may be set to an empty string ("").
	Name() string

	// ID returns the unique identifier of the mission.
	ID() string

	// GNGSetting returns the Go/No-Go setting for the mission.
	GNGSetting() GNGSetting
}

// Interface is the interface representing a falcon9 mission. This allows consumers
// to write their own mission logic if they wish to do so.
type Interface interface {
	InterfaceManageCrew
	InterfaceAccessors
	InterfaceLaunchControl
}

// MissionParams is a struct that consists of the parameters for a mission.
type MissionParams struct {
	ID     string
	GoNoGo GNGSetting
	Name   string
}

// Mission is the struct that implements the f9mission.Interface interface. This
// represents a falcon9 mission and all of its parameters.
type Mission struct {
	id   string
	name string
	gng  GNGSetting

	crew   map[string]f9crew.Interface
	crewMu sync.Mutex

	gngInProgress bool
	gngResults    Results
	gngMu         sync.Mutex
}

// NewMission is a function to provide a new mission with the provided parameters.
func NewMission(mp *MissionParams) (*Mission, error) {
	if mp == nil {
		return nil, errors.New("mission parameters cannot be nil")
	}

	if mp.ID == "" {
		return nil, errors.New("the ID of the mission cannot be an empty string")
	}

	m := &Mission{
		id:   mp.ID,
		name: mp.Name,
		gng:  mp.GoNoGo,
		crew: make(map[string]f9crew.Interface),
	}

	return m, nil
}

// Name returns the name of the mission. This is an optional value,
// so it may be set to an empty string ("").
func (m *Mission) Name() string { return m.name }

// ID returns the unique identifier of the mission.
func (m *Mission) ID() string { return m.id }

// GNGSetting returns the Go/No-Go setting for the mission.
func (m *Mission) GNGSetting() GNGSetting { return m.gng }

// Crew is a function that returns an f9crew.Manifest. This is a
// representation of the crew for the current mission. The slice
// returned contains unsorted values, but the returns value will
// have a Sort() method.
func (m *Mission) Crew() f9crew.Manifest {
	m.crewMu.Lock()
	defer m.crewMu.Unlock()

	manifest := make(f9crew.Manifest, len(m.crew))

	var counter int

	// loop over each item in the map
	// and append to the f9crew.Manifest slice
	for _, crew := range m.crew {
		manifest[counter] = crew
		counter++
	}

	return manifest
}

// AddCrew is a function to add a new crew member to this mission. If the crew
// member already exists (identified by their HashedKey), and replace is set to
// false, this will return an f9crew.ErrCrewMemberAlreadyPresent error. However,
// if replace is true it will simply replace the crew member with the new one.
//
// The latter is useful for a client quickly rejoining the session
// after a network interruption. This assume clients have a unique key.
func (m *Mission) AddCrew(crew f9crew.Interface, replace bool) error {
	// do some sanity checks before taking the mutex
	// if the crew map is nil, this struct was improperly created
	if m.crew == nil {
		return errUseNewMission
	}

	if crew == nil {
		return errors.New("a crew member cannot be nil")
	}

	m.crewMu.Lock()
	defer m.crewMu.Unlock()

	if _, ok := m.crew[crew.HashedKey()]; ok {
		// if we aren't going to replace the user
		// return an error
		if !replace {
			return ErrCrewMemberAlreadyPresent
		}

		delete(m.crew, crew.HashedKey())
	}

	m.crew[crew.HashedKey()] = crew

	return nil
}

// RemoveCrew is function to remove a crew member from the mission.
// The crew member's HashedKey is used to do the lookup for determining which
// crew member to remove from the mission. This returns the crew member being
// removed, if a consumer wishes to use it.
func (m *Mission) RemoveCrew(hashedKey string) (f9crew.Interface, error) {
	if hashedKey == "" {
		return nil, errors.New("hashedKey parameter cannot be an empty string")
	}

	m.crewMu.Lock()
	defer m.crewMu.Unlock()

	crew, ok := m.crew[hashedKey]

	if !ok {
		return nil, ErrCrewMemberNotPresent
	}

	delete(m.crew, hashedKey)

	return crew, nil
}

// Initiate is the function that starts the Go/No-Go call. At this point people
// can start adding votes to the mission.
func (m *Mission) Initiate() error {
	m.crewMu.Lock()

	n := len(m.crew)

	m.crewMu.Unlock()

	if n == 0 {
		return ErrNoAssignedCrew
	}

	m.gngMu.Lock()
	defer m.gngMu.Unlock()

	if m.gngInProgress {
		return ErrGNGInProgress
	}

	m.gngResults = make(Results)
	m.gngInProgress = true

	return nil
}

// UpdateVote updates the vote of a crew member for the current mission.
// The bool value returned indicates whether there have been enough "Go"
// votes to proceed with blastoff.
//
// If the mission is not initialized this will return a ErrGNGNotInProgress
// error. If the crew member is not assigned to this mission, this will return
// a ErrCrewMembeverNotPresent error.
func (m *Mission) UpdateVote(hashedKey string, vote Vote) (bool, error) {
	m.gngMu.Lock()
	defer m.gngMu.Unlock()

	if !m.gngInProgress {
		return false, ErrGNGNotInProgress
	}

	m.crewMu.Lock()
	defer m.crewMu.Unlock()

	if _, ok := m.crew[hashedKey]; !ok {
		return false, ErrCrewMemberNotPresent
	}

	m.gngResults[hashedKey] = vote

	return m.isReady(m.tally()), nil
}

func (m *Mission) isReady(t Tally) bool {
	numCrew := len(m.crew)

	switch m.gng {
	case GNGAll:
		return t[VoteYes] == numCrew
	case GNGQuorum:
		quroum := (numCrew / 2) + 1
		return t[VoteYes] >= quroum
	default:
		return false
	}
}

func (m *Mission) tally() Tally {
	tally := make(Tally)

	for _, vote := range m.gngResults {
		tally[vote]++
	}

	return tally
}

// Tally returns the tally of votes and whether there are enough votes
// to proceed with the mission.
func (m *Mission) Tally() (Tally, bool) {
	if !m.gngInProgress {
		return nil, false
	}

	m.gngMu.Lock()
	defer m.gngMu.Unlock()

	m.crewMu.Lock()
	defer m.crewMu.Unlock()

	tally := m.tally()
	return tally, m.isReady(tally)
}
