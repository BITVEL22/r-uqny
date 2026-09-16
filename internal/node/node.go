package node

// Node represents a r/uqny network participant.
type Node struct {
    ID string
}

// New creates a node with the given ID.
func New(id string) Node {
    return Node{
        ID: id,
    }
}
