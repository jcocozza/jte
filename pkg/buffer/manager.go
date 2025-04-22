package buffer

import (
	"fmt"
	"log/slog"
)

// A circular, doubly linked list
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
		n = newBufNode
		return n
	}
	last := n.prev
	newBufNode.next = n
	newBufNode.prev = last
	last.next = newBufNode
	n.prev = newBufNode
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

// keeps track of all buffers and maintains a pointer to the current (active) buffer
//
// the manager has 2 ways of keeping track of buffers
//  1. the bufList
//  2. the bufMap
//
// every operation should be safe for both
type BufferManager struct {
	Current *BufferNode
	bufList *BufferNode
	bufMap  map[int]*BufferNode
	// 'global' id counter
	// incremented each time a buffer is added
	// not decremented for any reason to keep id's unique
	idCounter int
	logger    *slog.Logger
}

func NewBufferManager(l *slog.Logger) *BufferManager {
	return &BufferManager{
		Current: nil,
		bufList: nil,
		bufMap: make(map[int]*BufferNode),
		logger: l.WithGroup("buffer-manager"),
	}
}

func (m *BufferManager) Add(buf *Buffer) int {
	m.idCounter++
	buf.id = m.idCounter
	newBufNode := m.bufList.Insert(buf)
	m.bufMap[newBufNode.id]	= newBufNode
	m.logger.Debug("add buffer", slog.Int("id", newBufNode.id))
	return newBufNode.id
}

func (m *BufferManager) Delete(id int) {
	m.logger.Debug("delete buffer", slog.Int("id", id))
	m.bufList.Delete(id)
	delete(m.bufMap, id)
}

func (m *BufferManager) SetCurrent(id int) {
	m.logger.Debug("set current buffer", slog.Int("id", id))
	if curr, ok := m.bufMap[id]; ok {
		m.Current = curr
		return
	}
	msg := fmt.Sprintf("invalid buffer id: %d", id)
	panic(msg)
}

func (m *BufferManager) Next() {
	msg := fmt.Sprintf("to next: %d -> %d", m.Current.id, m.Current.next.id)
	m.logger.Debug(msg)
	m.Current = m.Current.next
}

func (m *BufferManager) Previous() {
	msg := fmt.Sprintf("to prev: %d -> %d", m.Current.id, m.Current.prev.id)
	m.logger.Debug(msg)
	m.Current= m.Current.prev
}

type BufListData struct {
	Id int
	BufName string	
}

func (b *BufListData) String() string {
	return fmt.Sprintf("%d %s", b.Id, b.BufName)
}


func (m *BufferManager) ListAll() []BufListData {
	l := []BufListData{}
	for id, buf := range m.bufMap {
		b := BufListData{
			Id: id,
			BufName: buf.Buf.Name,
		}	
		l = append(l, b)	
	}
	return l
}
