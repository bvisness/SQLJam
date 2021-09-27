package sqljam


type SqlNode struct {
	X, Y      int
	//PinsInput []*SqlNodePin
	//PinsOutput []*SqlNodePin
}

/*
func (node *SqlNode) AddInputPin(pin *SqlNodePin) {
	pin.Parent = node
	node.PinsInput = append(node.PinsInput, pin)
}

func (node *SqlNode) AddOutputPin(pin *SqlNodePin) {
	pin.Parent = node
	node.PinsInput = append(node.PinsOutput, pin)
}
*/