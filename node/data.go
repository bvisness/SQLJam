package node

type NodeData interface {
	_impl() // delete if we ever actually put a meaningful method here
}

type Table struct {
	NodeData
	Table string
}

type PickColumns struct {
	NodeData
	Cols map[string]bool
}

type Filter struct {
	NodeData
	Conditions []string // TODO: whatever data we actually need for our filter UI
	// for now it always does AND because I am testing
}
