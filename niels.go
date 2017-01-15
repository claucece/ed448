package ed448

import "fmt"

type nielsTable []*bigNumber
type scalarAdjustmentTable [scalarWords]word_t //check this size

// A precomputed table for fixed-base scalar multiplication.
// This uses a signed combs format.
type fixedTable struct {
	base             nielsTable
	scalarAdjustment scalarAdjustmentTable
}

//XXX SECURITY this lookup should be done in constant time as it is on the
// original code
func (table *fixedTable) decafLookup(t, idx uint) *bigNumber {
	nin := 1 << (t - 1)    // j is always 1 now
	in := table.base[nin:] //this is not constant time
	return in[idx].copy()
}

var precomputedBaseTable = &fixedTable{}

// this is previously represented as niels point.
func init() {
	t := [4]*bigNumber{ // this is 240
		//0
		&bigNumber{0xf7e79ded, 0xd06e556e, 0xeedcf7ff, 0x6ec9befc, 0x27ee3efb, 0x32d7f7ff, 0x03c3f9ff, 0x05ee7dff},
		//1
		&bigNumber{0xdcbbcffb, 0x1cebd5cb, 0x8afeebf7, 0x117b5fff, 0x4b2fff7e, 0xd4eb11ef, 0x0a3ebaf7, 0xed5f7b5f},
		//2
		&bigNumber{0x79ef7f67, 0x73b77ceb, 0x6a5edee5, 0xaefda06a, 0x58fc754d, 0x47fd47f7, 0xedffaf7f, 0xf1f5eebf},
		//3
		&bigNumber{0x78feb30f, 0x21df6ffb, 0xf4ffffed, 0x79579fbf, 0x1ff76f77, 0x4cbdc76b, 0x70fd1fb7, 0x30ff76d3},
	}

	// check the gen of this
	adjustment := scalarAdjustmentTable{
		0x529eec33,
		0x721cf5b5,
		0xc8e9c2ab,
		0x7a4cf635,
		0x44a725bf,
		0xeec492d9,
		0x0cd77058,
		0x00000002,
		0x00000000,
		0x00000000,
		0x00000000,
		0x00000000,
		0x00000000,
		0x00000000,
	}

	precomputedBaseTable.base = t[:]
	precomputedBaseTable.scalarAdjustment = adjustment

}

func debugPrecomputedBaseTable() {
	for i, ni := range precomputedBaseTable.base {
		fmt.Printf("table[%d] %s\n", i, ni)
	}
}
