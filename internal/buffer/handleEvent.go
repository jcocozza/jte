package buffer

func (b *Buffer) RunningEvent() bool {
	if b.em.current == nil {
		return false
	}
	return !b.em.current.complete
}

func (b *Buffer) AcceptChange(c Change) error {
	err := c.Apply(b)
	if err != nil {
		return err
	}
	b.em.AddChange(c)
	return nil
}

func (b *Buffer) StartAndAcceptChange(c Change, etype EventType) error {
	err := c.Apply(b)
	if err != nil {
		return err
	}
	if b.em.current == nil {
		err := b.em.StartEvent(etype)
		if err != nil {
			return err
		}
	}
	b.em.AddChange(c)
	return nil
}

func (b *Buffer) Commit() {
	b.em.Commit()
}

func (b *Buffer) StartEvent(etype EventType) {
	b.em.StartEvent(etype)
}
