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

func step(c, s, m *bigNumber, n int64) {
	s.decafMul(m, c)
	c.copyFrom(s)
	for i := int64(0); i < n; i++ {
		c.decafSqr(c)
	}
}

func (n *bigNumber) decafIsqrt(x *bigNumber) {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	c.decafSqr(x)
	step(c, b, x, 1)
	step(c, b, x, 3)
	step(c, a, b, 3)
	step(c, a, b, 9)
	step(c, b, a, 1)
	step(c, a, x, 18)
	step(c, a, b, 37)
	step(c, b, a, 37)
	step(c, b, a, 111)
	step(c, a, b, 1)
	step(c, b, x, 223)
	n.decafMul(a, c)
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

// Weak reduce mod p
func (n *bigNumber) decafWeakReduce() {
	n[Limbs/2] += n[Limbs-1] >> Radix
	for i := uint(0); i < Limbs; i++ {
		n[i] += n[(i-1)%Limbs] >> Radix
		n[(i-1)%Limbs] &= radixMask
	}
}

// Substract mod p
// n = x - y
func (n *bigNumber) decafSub(x, y *bigNumber) {
	for i := 0; i < Limbs; i++ {
		n[i] = x[i] - y[i] + 2*P[i]
	}
	n.decafWeakReduce()
}

// Canonicalize, similar to strong reduce
func (n *bigNumber) decafCanon() {
	n.decafWeakReduce()

	carry := int64(0)

	for i := 0; i < Limbs; i++ {
		carry += int64(n[i]) - int64(P[i])
		n[i] = limb(carry) & radixMask
		carry >>= Radix
	}

	addback := carry
	carry = 0

	for j := 0; j < Limbs; j++ {
		carry += int64(n[j]) + (int64(P[j]) & addback)
		n[j] = limb(carry) & radixMask
		carry >>= Radix
	}
}

// compare x == y
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

// is neg? x : y
func (n *bigNumber) decafCondSel(x, y *bigNumber, neg word_t) {
	for i := uint(0); i < Limbs; i++ {
		n[i] = (x[i] & limb(^neg)) | (y[i] & limb(neg))
	}
}

func (n *bigNumber) decafCondNegate(neg word_t) {
	y := &bigNumber{}
	y.decafSub(&bigNumber{0}, n)
	n.decafCondSel(n, y, neg)
}
