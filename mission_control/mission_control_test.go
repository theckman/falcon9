package f9missioncontrol_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	seeded bool
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func seedRand(c *C) {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	fmt.Printf("random seed: %d\n", seed)
}

func (t *TestSuite) SetUpSuite(c *C) {
	if !t.seeded {
		seedRand(c)
		t.seeded = true
	}
}
