package hypercore

type Bitfield struct {
	pageSize int
	pager    Pager
}

func NewBitfield() Bitfield {
	bitfield := Bitfield{pageSize: 1024, pager: NewPager(1024)}
	return bitfield
}

// This is not correct at all.
func (b *Bitfield) Set(index int, value byte) {
	page_num := index / b.pageSize
	offset := index % b.pageSize
	page := b.pager.Get(page_num)
	buf := page.Buffer()
	buf[offset] = value
}

func (b *Bitfield) Get(index int) bool {
	return false
}
