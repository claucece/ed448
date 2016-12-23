package ed448

const (
	// D is the non-square element of F_p
	d                = -39081
	montgomeryFactor = "3bd440fae918bc5ull"
)

// P is biggish num
var P = []limb{radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask - 1, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask}

// working
// Copies n = y
func (n *bigNumber) decafCopy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

//failing: consts?
// XXX refactor, compare with Karatzuba mul, document what this does
// n = x * y
func (n *bigNumber) decafMul(x, y *bigNumber) {

	// copy so x is not directly modified
	xx := x.decafCopy()

	c := make([]dword_t, Limbs)

	for i := uint(0); i < Limbs; i++ {

		for j := uint(0); j < Limbs; j++ {
			c[(i+j)%Limbs] += dword_t(y[i]) * dword_t(xx[j])
		} // multiply and assigning in one value
		xx[(Limbs-1-i)^(Limbs/2)] += xx[Limbs-1-i] // assigning zeros
	}

	// shifting for mul
	c[Limbs-1] += c[Limbs-2] >> Radix
	c[Limbs-2] &= dword_t(radixMask) // masked off
	c[Limbs/2] += c[Limbs-1] >> Radix

	for k := uint(0); k < Limbs; k++ {
		c[k] += c[(k-1)%Limbs] >> Radix
		c[(k-1)%Limbs] &= dword_t(radixMask)
	}

	for l := uint(0); l < Limbs; l++ {
		n[l] = limb(c[l])
	}
}

func (n *bigNumber) decafMulW(x *bigNumber, y int64) {
	if y > 0 {
		yy := &bigNumber{limb(y)}
		n.decafMul(x, yy)
	} else {
		zz := &bigNumber{limb(-y)}
		n.decafMul(x, zz)
		zero := &bigNumber{0}
		n.decafSub(zero, n)
	}
}

// XXX this should be just twice mul
func (n *bigNumber) decafSqr(x *bigNumber, y uint) *bigNumber {
	sqr := n.squareN(x, y)
	return sqr
}

// working
// Weak reduce mod p
func (n *bigNumber) decafWeakReduce() {
	n[Limbs/2] += n[Limbs-1] >> Radix
	for i := uint(0); i < Limbs; i++ {
		n[i] += n[(i-1)%Limbs] >> Radix
		n[(i-1)%Limbs] &= radixMask
	}
}

// working
// Substract mod p
// n = x = y
func (n *bigNumber) decafSub(x, y *bigNumber) {
	for i := 0; i < Limbs; i++ {
		n[i] = x[i] - y[i] + 2*P[i] // maybe similar to bias?
	}
	n.decafWeakReduce()
}

// XXX: check if working complety
// Canonicalize, similar to strong reduce
func (n *bigNumber) decafCanon() {
	n.decafWeakReduce()

	// substract p with borrow

	// XXX: use d_word_t instead of word_t?
	var carry limb

	// XXX: guess there is no need for uint.. for range?
	for i := uint(0); i < Limbs; i++ {
		carry = carry + n[i] - P[i]
		n[i] &= carry & radixMask
		carry >>= Radix
	}

	addback := carry
	carry = 0

	// add it back
	for j := uint(0); j < Limbs; j++ {
		carry += carry + n[j] + (P[j] & addback)
		n[j] = carry & radixMask
		carry >>= Radix
	}
}

//working
// compare a == b
func decafEq(x, y *bigNumber) bool {
	n := &bigNumber{}
	n.decafSub(x, y)
	n.decafCanon()

	var ret limb

	for i := 0; i < Limbs; i++ {
		ret |= n[i]
	}
	return ((dword_t(ret) - 1) >> 32) != 0
}

func (n *bigNumber) decafAdd(x, y *bigNumber) {
	for i := uint(0); i < Limbs; i++ {
		n[i] = x[i] + y[i]
	}
	n.decafWeakReduce()
}
