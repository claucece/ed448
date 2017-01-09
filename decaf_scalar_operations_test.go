package ed448

import (
	. "gopkg.in/check.v1"
)

// XXX: check decaf_encode_from_ec in sage
func (s *Ed448Suite) Test_ValidateBasePoint(c *C) {

	c.Skip("this is the equation for ec25519")

	// This point representation is named extended
	// twisted Edwards coordinates and is denoted
	// by E^e. The identity element is represented by (0: 1: 0: 1).
	// The negative of (X: Y : T : Z) is (−X: Y : −T : Z).
	// Given (X: Y : Z) in E passing to E^e can be
	// performed in 3M+ 1S by computing (XZ, Y Z,XY,Z2). Given (X: Y : T : Z)
	// in E^e passing to E is cost-free by simply ignoring T .

	px := &bigNumber{
		0x00ffffff, 0xfffffffe,
		0x00ffffff, 0xffffffff,
		0x00ffffff, 0xffffffff,
		0x00ffffff, 0xffffffff,
		0x00000000, 0x00000003,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
	}

	py := &bigNumber{
		0x0081e6d3, 0x7f752992,
		0x003078ea, 0xd1c28721,
		0x00135cfd, 0x2394666c,
		0x0041149c, 0x50506061,
		0x0031d30e, 0x4f5490b3,
		0x00902014, 0x990dc141,
		0x0052341b, 0x04c1e328,
		0x00142378, 0x53c10a1b,
	}

	pz := &bigNumber{
		0x00ffffff, 0xfffffffb,
		0x00ffffff, 0xffffffff,
		0x00ffffff, 0xffffffff,
		0x00ffffff, 0xffffffff,
		0x00ffffff, 0xfffffffe,
		0x00ffffff, 0xffffffff,
		0x00ffffff, 0xffffffff,
		0x00ffffff, 0xffffffff,
	}

	pt := &bigNumber{
		0x008f205b, 0x70660415,
		0x00881c60, 0xcfd3824f,
		0x00377a63, 0x8d08500d,
		0x008c66d5, 0xd4672615,
		0x00e52fa5, 0x58e08e13,
		0x0087770a, 0xe1b6983d,
		0x004388f5, 0x5a0aa7ff,
		0x00b4d9a7, 0x85cf1a91,
	}

	p := &pointT{
		px,
		py,
		pz,
		pt,
	}

	// Z = F.random_element()
	// T = X*Y*Z
	// X = X*Z
	// Y = Y*Z
	// a := -1 for twisted .. not sure
	// a*X^2 + Y^2 == Z^2 + d*T^2
	r := &bigNumber{}
	t := &bigNumber{}
	r1 := &bigNumber{}
	a := &bigNumber{}
	r2 := &bigNumber{}
	r3 := &bigNumber{}
	r4 := &bigNumber{}
	a1 := &bigNumber{}
	b := int64(-1)
	w := &bigNumber{}
	w1 := &bigNumber{}

	r.decafSqr(px)
	t.decafMulW(r, b)
	r1.decafSqr(py)
	a.decafAdd(t, r1)

	r2.decafSqr(pz)
	r3.decafSqr(pt)
	r4.decafMulW(r3, -39081)
	a1.decafAdd(r2, r4)

	// y * pz = py because of projective coordinates
	// the Cartesian point (1, 2) can be represented
	// in homogeneous coordinates as (1, 2, 1) or (2, 4, 2).
	// The original Cartesian coordinates are recovered
	// by dividing the first two positions by the third.

	w.decafMul(px, py)
	w1.decafMul(pz, pt)

	c.Assert(decafPointValidate(p), Equals, word_t(0xffffffff))
	c.Assert(w, DeepEquals, w1)
	c.Assert(a, DeepEquals, a1)
}
