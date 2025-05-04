package keyboard

type OrderedKeyList []Key

func (o OrderedKeyList) Collapse() string {
	return string(o)
}

func (o *OrderedKeyList) Append(k Key) {
	*o = append(*o, k)
}
