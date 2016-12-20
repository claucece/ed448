package ed448

const (
	// D is the non-square element of F_p
	d                = -39081
	montgomeryFactor = "3bd440fae918bc5ull"
)

// P is bigish num
var P = []limb{radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask - 1, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask}

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
	n[Limbs-1] += n[Limbs-2] >> Radix
	n[Limbs-2] &= dword_t(radixMask) // masked off
	n[Limbs/2] += n[Limbs-1] >> Radix

	for k := uint(0); k < Limbs; k++ {
		n[k] += n[(k-1)%Limbs] >> Radix
		n[(k-1)%Limbs] &= dword_t(radixMask)
	}

	for l := uint(0); l < Limbs; l++ {
		c[l] = limb(n[l])
	}

	return c
}

// XXX this should be just twice mul
func (n *bigNumber) decafSqr(x *bigNumber, y uint) *bigNumber {
	sqr := n.squareN(x, y)
	return sqr
}

// working
//weak reduce mod p
func decafWeakReduce(x *bigNumber) *bigNumber {
	x[Limbs/2] += x[Limbs-1] >> Radix
	for i := uint(0); i < Limbs; i++ {
		x[i] += x[(i-1)%Limbs] >> Radix
		x[(i-1)%Limbs] &= radixMask
	}
	return x
}

// Substract mod p.
func decafSub(n, x, y *bigNumber) *bigNumber {
	for i := uint(0); i < Limbs; i++ {
		n[i] = x[i] - y[i] + 2*P[i] // maybe similar to bias?
		decafWeakReduce(n)
	}
	return n
}

// canonicalize
// similar to strong reduce
func decafCanon(x *bigNumber) *bigNumber {
	decafWeakReduce(x)

	// substract p with borrow

	// use d_word_t instead of word_t?
	var carry limb

	// guess there is no need for uint.. for range?
	for i := uint(0); i < Limbs; i++ {
		carry = carry + x[i] - P[i]
		x[i] &= carry & radixMask
		carry >>= Radix
	}

	addback := carry
	carry = 0

	// add it back
	for j := uint(0); j < Limbs; j++ {
		carry += carry + x[j] + (P[j] & addback)
		x[j] = carry & radixMask
		carry >>= Radix
	}

	return x
}

// compare a == b
func decafEq(x, y *bigNumber) dword_t {
	n := &bigNumber{}
	decafSub(n, x, y)
	decafCanon(n)

	var ret limb

	for i := uint(0); i < Limbs; i++ {
		ret |= n[i]
	}

	return dword_t(ret-1) >> 32
}

// should be a point
// uint32 bool
// are we using bools?
//func pointEq(p, q *bigNumber) bool {
//	a := &bigNumber{}
//	b := &bigNumber{}

//	decafMul(a, p, q)
//	decafMul(b, p, q)

//	return decafEq(a, b)
//}
