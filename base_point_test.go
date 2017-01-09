package ed448

import (
	"encoding/hex"
	"fmt"

	. "gopkg.in/check.v1"
)

// XXX: check decaf_encode_from_ec in sage
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

	y1, _ := decafDeser64(serialized{0x13})
	x1, _ := hex.DecodeString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f")

	x := new(bigNumber64).setBytes64(x1)
	fmt.Println(x)
	fmt.Println("the y", y1)
	z, _ := decafDeser64(serialized{0x01})
	fmt.Println("the z", z)

	t1 := &bigNumber64{}
	t := &bigNumber64{}

	t1.decafMul64(x, y1)
	t.decafMul64(t1, pz)

	r := &bigNumber64{}
	r.decafMul64(px, py)
	r.decafMul64(pz, pz)

	dst := [56]byte{}
	decafSerialize64(dst[:], px)

	fmt.Println(dst[:])
	// pt * pz = xy
	// px * py

	// y * pz = py because of projective coordinates
	// the Cartesian point (1, 2) can be represented
	// in homogeneous coordinates as (1, 2, 1) or (2, 4, 2).
	// The original Cartesian coordinates are recovered
	// by dividing the first two positions by the third.

	// pt = xy/z

	c.Assert(t, DeepEquals, pt)
	//c.Assert(r2, DeepEquals, py)
	//c.Assert(r3, DeepEquals, r)
	//c.Assert(h, DeepEquals, g)
}
