package container

import "sync"

type ComparableNode interface {
    Less(interface{}) bool
    String() string
}

type PriorityHeap struct {
    nodes     []ComparableNode
    lock      sync.RWMutex
    cachedKey map[string]bool
}

func NewPriorityHeap() (ph *PriorityHeap) {
    return &PriorityHeap{
        nodes:     make([]ComparableNode, 0, 1024),
        cachedKey: make(map[string]bool, 1024),
    }
}

func (self *PriorityHeap) Len() int {
    self.lock.RLock()
    defer self.lock.RUnlock()
    return len(self.nodes)
}

func (self *PriorityHeap) Enqueue(node ComparableNode) {
    if node == nil {
        return
    }

    self.lock.Lock()
    defer self.lock.Unlock()
    nodeKey := node.String()
    if self.cachedKey[nodeKey] {
        return
    }

    self.cachedKey[nodeKey] = true
    self.nodes = append(self.nodes, node)
    pos := len(self.nodes) - 1
    for {
        if pos == 0 {
            return
        }

        parentPos := pos / 2
        if !self.nodes[pos].Less(self.nodes[parentPos]) {
            return
        }

        tmpNode := self.nodes[parentPos]
        self.nodes[parentPos] = self.nodes[pos]
        self.nodes[pos] = tmpNode
        pos = parentPos
    }
}

func (self *PriorityHeap) Dequeue() (minNode ComparableNode) {
    self.lock.Lock()
    defer self.lock.Unlock()

    if len(self.nodes) == 0 {
        return
    }

    defer func() {
        nodeKey := minNode.String()
        delete(self.cachedKey, nodeKey)
    }()
    minNode = self.nodes[0]
    if len(self.nodes) == 1 {
        self.nodes = self.nodes[1:]
        return
    }

    nodeNum := len(self.nodes)
    self.nodes[0] = self.nodes[nodeNum-1]
    self.nodes = self.nodes[:nodeNum-1]
    nodeNum--
    rootPos := 0
    for {
        rootNode := self.nodes[rootPos]
        var leftSubNode, rightSubNode ComparableNode
        leftSubNodePos := (rootPos+1)*2 - 1
        if leftSubNodePos < nodeNum {
            leftSubNode = self.nodes[leftSubNodePos]
        }
        rightSubNodePos := (rootPos + 1) * 2
        if rightSubNodePos < nodeNum {
            rightSubNode = self.nodes[rightSubNodePos]
        }

        if leftSubNode == nil {
            return
        }

        minSubNode := leftSubNode
        minSubNodePos := leftSubNodePos
        if rightSubNode != nil && rightSubNode.Less(leftSubNode) {
            minSubNode = rightSubNode
            minSubNodePos = rightSubNodePos
        }

        // root 节点已经是最小的节点
        if rootNode.Less(minSubNode) {
            return
        }

        self.nodes[rootPos] = minSubNode
        self.nodes[minSubNodePos] = rootNode
        rootPos = minSubNodePos
    }
}
