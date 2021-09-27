package sqljam

import "log"

type SqlNodePin struct {
	Parent *SqlNode
	Connection *SqlNodePin
}

func (pin *SqlNodePin) ConnectTo(other *SqlNodePin) {
	if pin.Parent == nil || pin.Parent == nil {
		log.Panicf("We can't connect pins if one of the pins is parentless!")
	}
	pin.Connection = other
}

func (pin *SqlNodePin) Disconnect() {
	pin.Connection = nil
}
