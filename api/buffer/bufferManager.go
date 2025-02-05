package buffer

import "log/slog"

type BufNode struct {
	id   int
	Buf  Buffer
	next *BufNode
	prev *BufNode
}

// A circular doubly linked list
type BufList struct {
	node *BufNode
}

func (bl *BufList) Insert(buf Buffer, id int) *BufNode {
	newBufNode := &BufNode{
		id:   id,
		Buf:  buf,
		next: nil,
		prev: nil,
	}
	if bl.node == nil {
		newBufNode.next = newBufNode
		newBufNode.prev = newBufNode
		bl.node = newBufNode
		return newBufNode
	}
	last := bl.node.prev      // save the last node
	newBufNode.next = bl.node // point back to beginning
	newBufNode.prev = last    // newBufNode is the last, so it needs to point to what was the last node
	last.next = newBufNode    // old last now points to the new end
	bl.node.prev = newBufNode // the previous of the beginning circles to the end
	return newBufNode
}

func (bl *BufList) TraverseTo(id int) *BufNode {
	if bl.node == nil {
		return nil
	}
	curr := bl.node
	for curr.id != id {
		if curr.next == bl.node { // full circle without finding anything
			return nil
		}
		curr = bl.node.next
	}
	return curr
}

func (bl *BufList) Delete(id int) {
	if bl.node == nil {
		return
	}
	nodeToDelete := bl.node
	for nodeToDelete.id != id {
		if nodeToDelete.next == bl.node { // full circle without finding the node
			return
		}
		nodeToDelete = nodeToDelete.next
	}
	prev := nodeToDelete.prev
	next := nodeToDelete.next
	prev.next = next
	next.prev = prev
	if nodeToDelete == bl.node { // If node to delete is the first node
		bl.node = next
	}
	nodeToDelete.Buf.Close()
	nodeToDelete = nil
}

type BufferManager struct {
	bufferList  BufList
	bufMap      map[int]*BufNode
	CurrBufNode *BufNode
	// 'global' id counter.
	// incremented every time a buffer is added
	// not decremented for any reason as to keep id's unique
	idCounter int
	logger *slog.Logger
}

func NewBufferManager(l *slog.Logger) *BufferManager {
	bufList := BufList{node: nil}
	bufMap := make(map[int]*BufNode)
	return &BufferManager{
		bufferList: bufList,
		bufMap:     bufMap,
		logger: l.WithGroup("buffer manager"),
	}
}

func (bm *BufferManager) Add(buf Buffer) int {
	bm.idCounter++
	newBufNode := bm.bufferList.Insert(buf, bm.idCounter)
	bm.bufMap[newBufNode.id] = newBufNode
	return newBufNode.id
}

func (bm *BufferManager) SetCurrent(id int) {
	curr, ok := bm.bufMap[id]
	if !ok {
		return
	}
	bm.CurrBufNode = curr
}

func (bm *BufferManager) Delete(id int) {
	bm.bufferList.Delete(id)
	delete(bm.bufMap, id)
}

func (bm *BufferManager) Next() {
	bm.logger.Debug("going to next", slog.String("current", bm.CurrBufNode.Buf.Name()), slog.String("next", bm.CurrBufNode.next.Buf.Name()))
	bm.CurrBufNode = bm.CurrBufNode.next
}

func (bm *BufferManager) Prev() {
	bm.logger.Debug("going to prev", slog.String("current", bm.CurrBufNode.Buf.Name()), slog.String("previous", bm.CurrBufNode.prev.Buf.Name()))
	bm.CurrBufNode = bm.CurrBufNode.prev
}

type BufInfo struct {
	ID   int
	Name string
}

func (bm *BufferManager) List() []BufInfo {
	l := []BufInfo{}
	for id, bufn := range bm.bufMap {
		l = append(l, BufInfo{ID: id, Name: bufn.Buf.Name()})
	}
	return l
}
