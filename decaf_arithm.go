package ed448

const (
	//SerBytes is the number of bytes for serialization
	SerBytes = 56
)

//XXX: check var names and probably stop using word_t
//XXX: check overall order

//P is a biggish number
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

	for j := uint(0); j < Limbs; j++ {
		c[j] += c[(j-1)%Limbs] >> Radix
		c[(j-1)%Limbs] &= dword_t(radixMask)
	}

	for k := uint(0); k < Limbs; k++ {
		n[k] = limb(c[k])
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

// Substract with bias
// n = x - y
func (n *bigNumber) decafSubRawWithX(x, y *bigNumber, amt uint32) *bigNumber {
	return n.subRaw(x, y).bias(amt).weakReduce()
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

// XXX: this is one of the 3 func that decaf changes
// move it to proper place
// XXX: check the return value
// compare x == y
func decafEq(x, y *bigNumber) dword_t {
	n := &bigNumber{}
	n.decafSub(x, y)
	n.decafCanon()

	var ret limb

	for i := 0; i < Limbs; i++ {
		ret |= n[i]
	}
	return ((dword_t(ret) - 1) >> 32)
}

func (n *bigNumber) decafAdd(x, y *bigNumber) {
	for i := uint(0); i < Limbs; i++ {
		n[i] = x[i] + y[i]
	}
	n.decafWeakReduce()
}

func (n *bigNumber) decafCondSwap(x *bigNumber, swap dword_t) {
	for i := uint(0); i < Limbs; i++ {
		s := n[i] ^ x[i]&limb(swap)
		n[i] ^= s
		x[i] ^= s
	}

}

// is neg? x : y
func (n *bigNumber) decafCondNegate(neg dword_t) {
	y := &bigNumber{}
	y.sub(&bigNumber{0}, n)
	n.decafCondSel(n, y, neg)
}

func (n *bigNumber) decafCondSel(x, y *bigNumber, neg dword_t) {
	n[0] = (x[0] & limb(^neg)) | (y[0] & limb(neg))
	n[1] = (x[1] & limb(^neg)) | (y[1] & limb(neg))
	n[2] = (x[2] & limb(^neg)) | (y[2] & limb(neg))
	n[3] = (x[3] & limb(^neg)) | (y[3] & limb(neg))
	n[4] = (x[4] & limb(^neg)) | (y[4] & limb(neg))
	n[5] = (x[5] & limb(^neg)) | (y[5] & limb(neg))
	n[6] = (x[6] & limb(^neg)) | (y[6] & limb(neg))
	n[7] = (x[7] & limb(^neg)) | (y[7] & limb(neg))
	n[8] = (x[8] & limb(^neg)) | (y[8] & limb(neg))
	n[9] = (x[9] & limb(^neg)) | (y[9] & limb(neg))
	n[10] = (x[10] & limb(^neg)) | (y[10] & limb(neg))
	n[11] = (x[11] & limb(^neg)) | (y[11] & limb(neg))
	n[12] = (x[12] & limb(^neg)) | (y[12] & limb(neg))
	n[13] = (x[13] & limb(^neg)) | (y[13] & limb(neg))
	n[14] = (x[14] & limb(^neg)) | (y[14] & limb(neg))
	n[15] = (x[15] & limb(^neg)) | (y[15] & limb(neg))
}

// XXX: move this to proper place and probably divide it in two functions
// deserialize
// XXX: this is returning an extra 0
func decafDeser(in serialized) (*bigNumber, dword_t) {
	n := &bigNumber{}

	var k, bits uint
	var buf dword_t
	var accum dword_t

	for i := uint(0); i < SerBytes; i++ {
		buf |= dword_t(in[i]) << bits
		for bits += 8; (bits >= Radix || i == SerBytes-1) && k < Limbs; k, bits, buf = k+1, bits-Radix, buf>>Radix {
			n[k] = limb(buf) & radixMask
		}
	}

	for i := uint(0); i < Limbs; i++ {
		accum += (dword_t(n[i]) - dword_t(P[i])) >> wordBits
	}

	return n, accum
}
