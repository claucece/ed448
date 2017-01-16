package ed448

import "fmt"

type nielsTable []*twNiels
type scalarAdjustmentTable [scalarWords]word_t //check this size

// A precomputed table for fixed-base scalar multiplication.
// This uses a signed combs format.
type fixedTable struct {
	base             nielsTable
	scalarAdjustment scalarAdjustmentTable
}

//XXX SECURITY this lookup should be done in constant time as it is on the
// original code
func (table *fixedTable) decafLookup(j, t, idx uint) *twNiels {
	nin := 0 // for testing
	//nin := 1 << (t - 1)    // j is always 1 now
	in := table.base[nin:] //this is not constant time
	return in[idx].copy()
}

var precomputedBaseTable = &fixedTable{}

// is this the same if running the generation script with diff comb constants?
func init() {
	t := [3]*twNiels{ // 80.. is this correct? DECAF_COMBS_N<<(DECAF_COMBS_T-1)
		//0
		&twNiels{
			&bigNumber{0x07278dc5, 0x0e614a9f,
				0x004c5124, 0x02e454ad,
				0x0e1436f3, 0x0d8f58ce,
				0x0e4180ec, 0x0c83ed46,
				0x074a38fa, 0x0a41e932,
				0x0257771e, 0x0c1e7e53,
				0x03c0392f, 0x043e0ff0,
				0x05ce61df, 0x02c7c640,
			},
			&bigNumber{0x0c990b33, 0x033c4f9d,
				0x0ceb55c3, 0x0c291cb1,
				0x0ade88b2, 0x02ae3f58,
				0x01395474, 0x06b1f9f1,
				0x0b27ff7c, 0x02ded6e4,
				0x04aa10e1, 0x041012ed,
				0x0a36bae7, 0x03c22d20,
				0x0d472b19, 0x01f584ee,
			},
			&bigNumber{0x09ee6f60, 0x0c351477,
				0x03b20c2b, 0x01574c87,
				0x0a5a5e65, 0x04cd6a46,
				0x0eb4204a, 0x059a068a,
				0x08bc354d, 0x04c61045,
				0x079d02d2, 0x0e945674,
				0x0d118e28, 0x0feaf77e,
				0x0115eeb5, 0x0f58a8bf,
			},
		},
		//9
		&twNiels{
			&bigNumber{0x0ad825f1, 0x0d37716c,
				0x0ba9552a, 0x0883870c,
				0x05c762e3, 0x08ef785f,
				0x00469242, 0x06cb253e,
				0x0ee9d967, 0x07b8f17f,
				0x032b52b6, 0x0a43de69,
				0x02af783c, 0x01aca9fe,
				0x0ff0b680, 0x08967778,
			},
			&bigNumber{0x0dc6c9c3, 0x06400c4c,
				0x0691083f, 0x01e8c978,
				0x0f68e0c5, 0x0ad74f01,
				0x072b5f6a, 0x0f7feb03,
				0x05ade13a, 0x02f60d17,
				0x0221a678, 0x098ec54a,
				0x071f244e, 0x0fcfea8a,
				0x0e45ded2, 0x0dea6660,
			},
			&bigNumber{0x0a8d6752, 0x02585b4a,
				0x015a2089, 0x0e62da76,
				0x01f39b68, 0x010c1c74,
				0x0ced9f65, 0x0569bb1e,
				0x04daa724, 0x0ba6d09e,
				0x0ef281b9, 0x07d3e20a,
				0x0ca3ffdc, 0x0bd7f65a,
				0x050288a8, 0x0dea434a,
			},
		},
		//16
		&twNiels{
			&bigNumber{0x084da36e, 0x03c97e63,
				0x0ac81a09, 0x0423d53e,
				0x03cdce35, 0x0b70d68f,
				0x0354b92c, 0x0ee7959b,
				0x0819c8ca, 0x0f4e9718,
				0x0acbffe9, 0x09349f12,
				0x02cb7da6, 0x05aee7b6,
				0x054ffc86, 0x0d977641,
			},
			&bigNumber{0x0fcb435a, 0x0d95d1c5,
				0x0b5086f9, 0x016d1ed6,
				0x07e54d71, 0x0792aa0b,
				0x05f1925d, 0x067b6571,
				0x0ec6176b, 0x0a219755,
				0x0b12c28f, 0x0bc3f026,
				0x0ffeb93e, 0x0700c897,
				0x0ec50b46, 0x089b83f6,
			},
			&bigNumber{0x0544b923, 0x0ad9cdb4,
				0x07284061, 0x0d11664c,
				0x0b8f910b, 0x0815ae86,
				0x0591c3c6, 0x05414fb2,
				0x02d7ef9e, 0x094ba83e,
				0x0599386c, 0x001dbc16,
				0x0493911b, 0x0c8721f0,
				0x063c346c, 0x0c1be6b4,
			},
		},
	}
	precomputedBaseTable.base = t[:]
}

func debugPrecomputedBaseTable() {
	for i, ni := range precomputedBaseTable.base {
		fmt.Printf("table[%d] %s\n", i, ni)
	}
}

func (n *twNiels) condNegNiels(neg word_t) {
	n.a.decafCondSwap(n.b, neg)
	n.c.decafCondNegate(neg)
}
