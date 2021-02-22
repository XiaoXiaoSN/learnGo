package main

import (
	"fmt"

	"github.com/looplab/fsm"
)

type Door struct {
	To       string
	Msg      string
	NewOrder *Door
	FSM      *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
			{Name: "end", Src: []string{"closed"}, Dst: "end"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				d.enterState(e, d.Msg)
				d.Msg = ""
			},
		},
	)

	return d
}

func (d *Door) enterState(e *fsm.Event, msg string) {
	fmt.Printf("[%s] The door to %s is %s\n", msg, d.To, e.Dst)
}

func main() {
	door := NewDoor("heaven")

	door.Msg = "open message"
	err := door.FSM.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	err = door.FSM.Event("close")
	if err != nil {
		fmt.Println(err)
	}

	door.FSM.SetState("end")
	fmt.Println("open >", door.FSM.Can("open"))
	err = door.FSM.Event("end")
	if err != nil {
		fmt.Println(err)
	}
}
