package buffer

// a circular, doubly linked list
//
// wraps the buffer and lets us navigate across a list of buffers
//
// in most cases, this is the point of interaction with the buffer for other parts of the code
// it allows access to the buffer and the next/previous
type BufferNode struct {
	id   int
	Buf  *Buffer
	next *BufferNode
	prev *BufferNode
}

func (n *BufferNode) Insert(buf *Buffer) *BufferNode {
	newBufNode := &BufferNode{
		id:   buf.id,
		Buf:  buf,
		next: nil,
		prev: nil,
	}
	if n == nil {
		newBufNode.next = newBufNode
		newBufNode.prev = newBufNode
		return newBufNode
	}
	next := n.prev

	newBufNode.prev = n
	newBufNode.next = next
	n.next = newBufNode
	next.prev = newBufNode
	return newBufNode
}

func (n *BufferNode) Delete(id int) *BufferNode {
	if n == nil {
		return nil
	}
	nodeToDelete := n
	for nodeToDelete.id != id {
		if nodeToDelete.next == n { // full circle without finding the node
			return n
		}
		nodeToDelete = nodeToDelete.next
	}
	prev := nodeToDelete.prev
	next := nodeToDelete.next
	prev.next = next
	next.prev = prev
	if nodeToDelete == n {
		return next
	}
	// unlink to make it clear to gc to clean up
	nodeToDelete.next = nil
	nodeToDelete.prev = nil
	return n
}

func (n *BufferNode) TraverseTo(id int) *BufferNode {
	if n == nil {
		return nil
	}
	curr := n
	for curr.id != id {
		if curr.next == n { // full circle without finding anything
			return nil
		}
		curr = n.next
	}
	return curr
}

func (n *BufferNode) Next() *BufferNode {
	return n.next
}

func (n *BufferNode) Previous() *BufferNode {
	return n.prev
}
