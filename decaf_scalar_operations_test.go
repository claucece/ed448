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

	gx, _ := decafDeser(serialized{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	})
	//gy := serialized{0x13}

	bigNumOne, _ := decafDeser(serialized{0x01})

	l := &bigNumber{}

	l.decafMul(gx, bigNumOne)

	c.Assert(decafPointValidate(p), Equals, word_t(0xffffffff))
	c.Assert(l, DeepEquals, px)
	c.Assert(w, DeepEquals, w1)
	c.Assert(a, DeepEquals, a1)
}

// this might work

//

//impl ProjectivePoint {
/// Convert to the extended twisted Edwards representation of this
/// point.
///
/// From §3 in [0]:
///
/// Given (X:Y:Z) in Ɛ, passing to Ɛₑ can be performed in 3M+1S by
/// computing (XZ,YZ,XY,Z²).  (Note that in that paper, points are
/// (X:Y:T:Z) so this really does match the code below).
//    #[allow(dead_code)]  // rustc complains this is unused even when it's used
//   fn to_extended(&self) -> ExtendedPoint {
//        ExtendedPoint{
//            X: &self.X * &self.Z,
//            Y: &self.Y * &self.Z,
//            Z: self.Z.square(),
//            T: &self.X * &self.Y,
//        }
//    }

func (s *Ed448Suite) Test_ScalarOperations(c *C) {

	scalar1 := [scalarWords]word_t{
		50, 0, 0, 0, 6, 0, 0, 3, 0, 0, 0, 2, 1, 1,
	}

	scalar2 := [scalarWords]word_t{
		5, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1,
	}

	subExp := [scalarWords]word_t{
		45, 0, 0, 0, 6, 0, 0, 1, 0, 0, 0, 2, 1, 0,
	}

	addExp := [scalarWords]word_t{
		55, 0, 0, 0, 6, 0, 0, 5, 0, 0, 0, 2, 1, 2,
	}

	added := scalarAdd(scalar1, scalar2)
	subtracted := scalarSub(scalar1, scalar2)

	c.Assert(added, DeepEquals, addExp)
	c.Assert(subtracted, DeepEquals, subExp)
}

func (s *Ed448Suite) Test_GenerateConstant(c *C) {

	c.Skip("In progress")
	//constant := [scalarWords]word_t{
	//	0x4a7bb0cf, 0xc873d6d5, 0x23a70aad, 0xe933d8d7, 0x129c96fd, 0xbb124b65, 0x335dc163,
	//	0x00000008, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
	//}

	//scalar := scalarAdjustment()

	//c.Assert(constant, DeepEquals, scalar)
}
