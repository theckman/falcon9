package f9missioncontrol

import (
	"fmt"
	"sync"

	"github.com/theckman/falcon9/mission"
)

type missionRegistry struct {
	missions map[uint16]f9mission.Interface
}

var (
	registry   missionRegistry
	registryMu sync.RWMutex
)

// GetMission returns a mission, based on the ID, if one has been created. If the
// mission doesn't exist this just returns nil.
func GetMission(id uint16) f9mission.Interface {
	registryMu.RLock()
	defer registryMu.RUnlock()

	ifc, ok := registry.missions[id]

	if !ok {
		return nil
	}

	return ifc
}

// AddMission is a function to add a mission to the registry. This will only
// return an error when the registry already has a mission with that ID.
func AddMission(id uint16, mission f9mission.Interface) error {
	registryMu.Lock()
	defer registryMu.Unlock()

	if _, ok := registry.missions[id]; ok {
		return fmt.Errorf("Mission with ID %d already registered", id)
	}

	registry.missions[id] = mission

	return nil
}

// RemoveMission purges a mission from the mission registry. If the mission existed
// this will return the mission, otherwise it will return nil.
func RemoveMission(id uint16) f9mission.Interface {
	registryMu.Lock()
	defer registryMu.Unlock()

	if mission, ok := registry.missions[id]; ok {
		delete(registry.missions, id)
		return mission
	}

	return nil
}

// ListMissions returns a slice of the mission IDs. They are in no particular order.
func ListMissions() []uint16 {
	registryMu.RLock()
	defer registryMu.RUnlock()

	slice := make([]uint16, len(registry.missions))

	var i int

	for id := range registry.missions {
		slice[i] = id
		i++
	}

	return slice
}

func init() {
	registry.missions = make(map[uint16]f9mission.Interface)
}
