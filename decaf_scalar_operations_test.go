package ed448

import . "gopkg.in/check.v1"

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
