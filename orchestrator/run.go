package orchestrator

import (
	"github.com/owulveryck/gorchestrator/structure"
	"math/rand"
	"time"
)

// Run executes the artifact of a given node
func (n *Node) Run() <-chan Message {
	c := make(chan Message)
	waitForIt := make(chan structure.Matrix) // Shared between all messages.
	go func() {
		state := ToRun
		for state <= ToRun {
			c <- Message{n.ID, state, waitForIt}
			m := <-waitForIt
			s := m.Dim()
			state = Running
			for i := 0; i < s; i++ {
				if m.At(i, n.ID) < Success && m.At(i, n.ID) > 0 {
					state = ToRun
				} else if m.At(i, n.ID) >= Failure {
					state = NotRunnable
					continue
				}
			}
			if state == NotRunnable {
				//fmt.Printf("I am %v, and I cannot run\n", n.ID)
				c <- Message{n.ID, state, waitForIt}
			}
			if state == Running {
				c <- Message{n.ID, state, waitForIt}
				//fmt.Printf("I am %v, and I am running: the module %v, with %v %v\n", n.ID, n.Engine, n.Artifact, n.Args)
				time.Sleep(time.Duration(rand.Intn(1e4)) * time.Millisecond)
				rand.Seed(time.Now().Unix())
				state = Success
				// Now send the message that I'm done...
				c <- Message{n.ID, state, waitForIt}
			}
		}
		close(c)
	}()
	return c
}
