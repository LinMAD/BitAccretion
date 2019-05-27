package dashboard

import (
	"github.com/LinMAD/BitAccretion/model"
	"testing"
)

func TestCreateLeftLayout(t *testing.T) {
	nodes := make([]model.Node, 4)
	nodes[0] = model.Node{Name: "System 0"}
	nodes[1] = model.Node{Name: "System 1"}
	nodes[2] = model.Node{Name: "System 2"}
	nodes[3] = model.Node{Name: "System 3"}

	left, e := CreateLeftLayout(nodes)
	if e != nil {
		t.Fail()
	}

	if left == nil {
		t.Fail()
	}
}