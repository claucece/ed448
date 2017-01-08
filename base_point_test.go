package ed448

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestBasePoint(c *C) {

	c.Skip("not sure")

	px := &bigNumber64{
		0x00fffffffffffffe, 0x00ffffffffffffff,
		0x00ffffffffffffff, 0x00ffffffffffffff,
		0x0000000000000003, 0x0000000000000000,
		0x0000000000000000, 0x0000000000000000,
	}

	py := &bigNumber64{
		0x0081e6d37f752992, 0x003078ead1c28721,
		0x00135cfd2394666c, 0x0041149c50506061,
		0x0031d30e4f5490b3, 0x00902014990dc141,
		0x0052341b04c1e328, 0x0014237853c10a1b,
	}

	pz := &bigNumber64{
		0x00fffffffffffffb, 0x00ffffffffffffff,
		0x00ffffffffffffff, 0x00ffffffffffffff,
		0x00fffffffffffffe, 0x00ffffffffffffff,
		0x00ffffffffffffff, 0x00ffffffffffffff,
	}

	pt := &bigNumber64{
		0x008f205b70660415, 0x00881c60cfd3824f,
		0x00377a638d08500d, 0x008c66d5d4672615,
		0x00e52fa558e08e13, 0x0087770ae1b6983d,
		0x004388f55a0aa7ff, 0x00b4d9a785cf1a91,
	}

	y := &bigNumber64{0x00000000000000013}

	r := &bigNumber64{}
	r2 := &bigNumber64{}
	r3 := &bigNumber64{}

	// pt * pz = xy
	// px * py
	r.decafMul64(px, py)

	// y * pz = py because of projective coordinates
	// the Cartesian point (1, 2) can be represented
	// in homogeneous coordinates as (1, 2, 1) or (2, 4, 2).
	// The original Cartesian coordinates are recovered
	// by dividing the first two positions by the third.
	r2.decafMul64(y, pz)

	// pt = xy/z
	r3.decafMul64(pt, pz)

	c.Assert(r, DeepEquals, pt)
	c.Assert(r2, DeepEquals, py)
	c.Assert(r3, DeepEquals, r)
}
