package ed448

// P is biggish num
var P = []limb{radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask - 1, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask, radixMask}

func (n *bigNumber) copyFrom(right *bigNumber) {
	copy(n[:], right[:])
}

// n = x * y
func (n *bigNumber) decafMul(x, y *bigNumber) {

	xx := x.copy()

	c := make([]dword_t, Limbs)

	for i := uint(0); i < Limbs; i++ {

		for j := uint(0); j < Limbs; j++ {
			c[(i+j)%Limbs] += dword_t(y[i]) * dword_t(xx[j])
		}
		xx[(Limbs-1-i)^(Limbs/2)] += xx[Limbs-1-i]
	}

	c[Limbs-1] += c[Limbs-2] >> Radix
	c[Limbs-2] &= dword_t(radixMask)
	c[Limbs/2] += c[Limbs-1] >> Radix

	for k := uint(0); k < Limbs; k++ {
		c[k] += c[(k-1)%Limbs] >> Radix
		c[(k-1)%Limbs] &= dword_t(radixMask)
	}

	for l := uint(0); l < Limbs; l++ {
		n[l] = limb(c[l])
	}
}

func (n *bigNumber) decafSqr(x *bigNumber) {
	n.decafMul(x, x)
}

func step(n, x, y *bigNumber, i int64) {
	x.decafMul(y, n)
	n = x.decafCopy()
	for j := int64(0); j < i; j++ {
		n.decafSqr(n)
	}
}

//func (n *bigNumber) decafIsr(y *bigNumber) {
//	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
//	c.decafSqr(y)
//	step(c, b, y, 1)
//	y.decafMul(a, c)
//}

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
