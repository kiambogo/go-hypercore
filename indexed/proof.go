package indexed

type Proof struct {
	index      uint64
	verifiedBy uint64
	nodes      []uint64
}
