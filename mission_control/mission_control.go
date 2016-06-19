package f9missioncontrol

import (
	"net"

	"github.com/theckman/falcon9/crew"
	"github.com/theckman/falcon9/mission"
)

type client struct {
	conn net.Conn
	out  chan []byte
	crew f9crew.Interface
}

// MissionControl is the controller of a mission.
type MissionControl struct {
	Mission f9mission.Interface
	clients map[string]*client
}
