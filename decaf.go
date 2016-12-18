package ed448

const (
	// D is the non-square element of F_p
	d                = -39081
	montgomeryFactor = "3bd440fae918bc5ull"
	lbits            = 28                     // bit field
	lmask            = (limb(1) << lbits) - 1 // masking to 0
)

// Copy copies n = y
func (n *bigNumber) decafCopy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

// XXX refactor, compare with Karatzuba mul, document what this does
func decafMul(n, x, y *bigNumber) *bigNumber {

	// copy so x is not directly modified
	xx := x.decafCopy()

	for i := 0; i < Limbs; i++ {
		for j := 0; j < Limbs; j++ {
			n[(i+j)%Limbs] += y[i] * xx[j]
			xx[(Limbs-1-i)^(Limbs/2)] += xx[Limbs-1-i]
		}
	}

	// shifting for mul
	n[Limbs-1] += n[Limbs-2] >> lbits
	n[Limbs-2] &= lmask // masked off
	n[Limbs/2] += n[Limbs-1] >> lbits

	// WHY?
	for k := 0; k < Limbs; k++ {
		if k != 0 {
			n[k] += n[(k-1)%Limbs] >> lbits
			n[(k-1)%Limbs] &= lmask
		} else {
			n[k] += n[(k)%Limbs] >> lbits
			n[(k)%Limbs] &= lmask
		}
	}

	return n
}

// XXX until we found out the new implementation of MH
func (n *bigNumber) decafSqr(x *bigNumber, y uint) *bigNumber {
	sqr := n.squareN(x, y)
	return sqr
}
