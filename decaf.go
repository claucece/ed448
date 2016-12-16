package ed448

const (
	// D is the non-square element of F_p
	d                = -39081
	montgomeryFactor = "3bd440fae918bc5ull"
	lbits            = 28
	lmask            = (limb(1) << lbits) - 1
)

// Copy copies n = y
func (n *bigNumber) decafCopy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

// XXX refactor, compare with Karatzuba mul, document what this does
func (n *bigNumber) decafMul(x, y *bigNumber) *bigNumber {

	xx := x.decafCopy()

	for i := 0; i < Limbs; i++ {
		for j := 0; j < Limbs; j++ {
			n[(i+j)%Limbs] += y[i] * xx[i]
			xx[(Limbs-1-i)^(Limbs/2)] += xx[Limbs-1-i]
		}
	}

	n[Limbs-1] += n[Limbs-2] >> lbits
	n[Limbs-2] &= lmask
	n[Limbs/2] += n[Limbs-1] >> lbits

	// WHY?
	for k := 1; k < Limbs; k++ {
		n[k] += n[(k-1)%Limbs] >> lbits
		n[(k-1)%Limbs] &= lmask
	}

	return n
}
