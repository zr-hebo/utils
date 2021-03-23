package container

import (
    "fmt"
    "testing"
)

type baseNode struct {
    val int
}

func (self *baseNode) Less(cmp interface{}) bool {
    if cmp == nil {
        return false
    }
    cbn := cmp.(*baseNode)
    return self.val < cbn.val
}

func (self *baseNode) String() string {
    return fmt.Sprint(self.val)
}

func TestForPriorityHeap(t *testing.T) {
    ph := NewPriorityHeap()
    vals := []int{1, 9, 15, 21, 27}
    for _, val := range vals {
        ph.Enqueue(&baseNode{
            val: val,
        })
    }
    node := ph.Dequeue()
    t.Log(node.String())

    vals = []int{30, 60, 90}
    for _, val := range vals {
        ph.Enqueue(&baseNode{
            val: val,
        })
    }
    node = ph.Dequeue()
    t.Log(node.String())

    vals = []int{24, 18, 12,6}
    for _, val := range vals {
        ph.Enqueue(&baseNode{
            val: val,
        })
    }

    bn := node.(*baseNode)
    minVal := 0
    orderVals := make([]int, 0, 16)
    for ph.Len() > 0 {
        node = ph.Dequeue()
        bn = node.(*baseNode)
        orderVals = append(orderVals, bn.val)
        if minVal > bn.val {
            t.Fatalf("val %v out of order, %v", bn.val, orderVals)
        }

        t.Log(node.String())
    }
}
