package ed448

import . "gopkg.in/check.v1"

type DecafNielsSuite struct{}

var _ = Suite(&DecafNielsSuite{})

func (s *DecafNielsSuite) Test_DecafLookup(c *C) {

	expA := &bigNumber{
		0x0ad825f1, 0x0d37716c,
		0x0ba9552a, 0x0883870c,
		0x05c762e3, 0x08ef785f,
		0x00469242, 0x06cb253e,
		0x0ee9d967, 0x07b8f17f,
		0x032b52b6, 0x0a43de69,
		0x02af783c, 0x01aca9fe,
		0x0ff0b680, 0x08967778,
	}

	expB := &bigNumber{0x0dc6c9c3, 0x06400c4c,
		0x0691083f, 0x01e8c978,
		0x0f68e0c5, 0x0ad74f01,
		0x072b5f6a, 0x0f7feb03,
		0x05ade13a, 0x02f60d17,
		0x0221a678, 0x098ec54a,
		0x071f244e, 0x0fcfea8a,
		0x0e45ded2, 0x0dea6660,
	}
	expC := &bigNumber{0x0a8d6752, 0x02585b4a,
		0x015a2089, 0x0e62da76,
		0x01f39b68, 0x010c1c74,
		0x0ced9f65, 0x0569bb1e,
		0x04daa724, 0x0ba6d09e,
		0x0ef281b9, 0x07d3e20a,
		0x0ca3ffdc, 0x0bd7f65a,
		0x050288a8, 0x0dea434a,
	}

	zerothNiels := uint(0)
	sixteenthNiels := uint(16)
	ninethNiels := uint(9)

	point := precomputedBaseTable.decafLookup(zerothNiels, sixteenthNiels, ninethNiels)

	c.Assert(expA, DeepEquals, point.a)
	c.Assert(expB, DeepEquals, point.b)
	c.Assert(expC, DeepEquals, point.c)

}

func (s *DecafNielsSuite) Test_DecafCondNegNiels(c *C) {

	n := &twNiels{
		&bigNumber{0x00000000},
		&bigNumber{0x0ac67eac, 0x08c3224f,
			0x038fe548, 0x09a46a59,
			0x0e30ed3f, 0x032c1eb2,
			0x08ebe610, 0x03168199,
			0x0dd4e788, 0x06d5a576,
			0x077ec52f, 0x00987f7d,
			0x03a54795, 0x08cbe066,
			0x0db4e599, 0x0af8126b,
		},
		&bigNumber{0x08db85c2, 0x0fd2361e,
			0x0ce2105d, 0x06a17729,
			0x0e3ca84d, 0x0a137aa5,
			0x0985ee61, 0x05a26d64,
			0x0734c5f3, 0x0da853af,
			0x01d955b7, 0x03160ecd,
			0x0a59046d, 0x0c32cf71,
			0x98dce72d, 0x00007fff,
		},
	}

	expA := &bigNumber{0x00000000}
	expB := &bigNumber{0x0ac67eac, 0x08c3224f,
		0x038fe548, 0x09a46a59,
		0x0e30ed3f, 0x032c1eb2,
		0x08ebe610, 0x03168199,
		0x0dd4e788, 0x06d5a576,
		0x077ec52f, 0x00987f7d,
		0x03a54795, 0x08cbe066,
		0x0db4e599, 0x0af8126b,
	}

	expC := &bigNumber{0x08db85c2, 0x0fd2361e,
		0x0ce2105d, 0x06a17729,
		0x0e3ca84d, 0x0a137aa5,
		0x0985ee61, 0x05a26d64,
		0x0734c5f3, 0x0da853af,
		0x01d955b7, 0x03160ecd,
		0x0a59046d, 0x0c32cf71,
		0x98dce72d, 0x00007fff,
	}

	n.condNegNiels(word_t(0))

	c.Assert(expA, DeepEquals, n.a)
	c.Assert(expB, DeepEquals, n.b)
	c.Assert(expC, DeepEquals, n.c)

}
