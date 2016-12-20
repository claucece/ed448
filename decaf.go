package ed448

const (
	// D is the non-square element of F_p
	d                = -39081
	montgomeryFactor = "3bd440fae918bc5ull"
	lbits            = 28                          // bit field
	lmask            = ((dword_t(1) << lbits) - 1) // masking to 0
)

// Copy copies n = y
func (n *bigNumber) decafCopy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

// XXX refactor, compare with Karatzuba mul, document what this does
// n = x * y
func decafMul(c, x, y *bigNumber) *bigNumber {

	// copy so x is not directly modified
	xx := x.decafCopy()

	n := make([]dword_t, Limbs)

	for i := uint(0); i < Limbs; i++ {
		for j := uint(0); j < Limbs; j++ {
			n[(i+j)%Limbs] += dword_t(y[i] * xx[j])
		} // multiply and assing in one value
		xx[(Limbs-1-i)^(Limbs/2)] += xx[Limbs-1-i] // assinging zeros
	}

	// shifting for mul
	n[Limbs-1] += n[Limbs-2] >> lbits
	n[Limbs-2] &= lmask // masked off
	n[Limbs/2] += n[Limbs-1] >> lbits

	for k := uint(0); k < Limbs; k++ {
		n[k] += n[(k-1)%Limbs] >> lbits
		n[(k-1)%Limbs] &= lmask
	}

	for l := uint(0); l < Limbs; l++ {
		c[l] = limb(n[l])
	}

	return c
}

// XXX until we found out the new implementation of MH
func (n *bigNumber) decafSqr(x *bigNumber, y uint) *bigNumber {
	sqr := n.squareN(x, y)
	return sqr
}
