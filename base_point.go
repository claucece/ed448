package ed448

// XXX: this probably need a setBytes if I want to play with the hex
// 0x297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f,
// 0x13
const (
	lmask = (dword_t(1) << 56) - 1
)

// P64 is P
var P64 = []dword_t{lmask, lmask, lmask, lmask, lmask - 1, lmask, lmask, lmask}

type bigNumber64 [8]dword_t

type base struct {
	x, y, z, t *bigNumber64
}

func (n *bigNumber64) copy64() *bigNumber64 {
	c := &bigNumber64{}
	copy(c[:], n[:])
	return c
}

func (n *bigNumber64) decafWeakReduce64() {
	n[8/2] += n[8-1] >> 56
	for i := uint(0); i < 8; i++ {
		n[i] += n[(i-1)%8] >> 56
		n[(i-1)%8] &= lmask
	}
}

func (n *bigNumber64) decafAdd64(x, y *bigNumber64) {
	for i := uint(0); i < 8; i++ {
		n[i] = x[i] + y[i]
	}
	n.decafWeakReduce64()
}

func (n *bigNumber64) decafSub64(x, y *bigNumber64) {
	for i := 0; i < 8; i++ {
		n[i] = x[i] - y[i] + 2*P64[i]
	}
	n.decafWeakReduce64()
}

func (n *bigNumber64) decafMulW64(x *bigNumber64, y int64) {
	if y > 0 {
		yy := &bigNumber64{dword_t(y)}
		n.decafMul64(x, yy)
	} else {
		zz := &bigNumber64{dword_t(-y)}
		n.decafMul64(x, zz)
		zero := &bigNumber64{0}
		n.decafSub64(zero, n)
	}
}

func (n *bigNumber64) decafMul64(x, y *bigNumber64) {
	xx := x.copy64()

	c := make([]dword_t, 8)

	for i := uint(0); i < 8; i++ {

		for j := uint(0); j < 8; j++ {
			c[(i+j)%8] += y[i] * xx[j]
		}
		xx[(8-1-i)^(8/2)] += xx[8-1-i]
	}

	c[8-1] += c[8-2] >> 56
	c[8-2] &= dword_t(lmask)
	c[8/2] += c[8-1] >> 56

	for j := uint(0); j < 8; j++ {
		c[j] += c[(j-1)%8] >> 56
		c[(j-1)%8] &= dword_t(lmask)
	}

	for k := uint(0); k < 8; k++ {
		n[k] = c[k]
	}
}

func (n *bigNumber64) decafCanon64() {
	n.decafWeakReduce64()

	carry := int64(0)

	for i := 0; i < 8; i++ {
		carry += int64(n[i]) - int64(P[i])
		n[i] = dword_t(carry) & lmask
		carry >>= 56
	}

	addback := carry
	carry = 0

	for j := 0; j < 8; j++ {
		carry += int64(n[j]) + (int64(P[j]) & addback)
		n[j] = dword_t(carry) & lmask
		carry >>= 56
	}
}

func decafDeser64(in serialized) (*bigNumber64, dword_t) {
	n := &bigNumber64{}

	var k, bits uint
	var buf dword_t
	var accum dword_t

	for i := uint(0); i < 56; i++ {
		buf |= dword_t(in[i]) << bits
		for bits += 8; (bits >= 56 || i == 56-1) && k < 8; k, bits, buf = k+1, bits-56, buf>>56 {
			n[k] = buf & lmask
		}
	}

	for i := uint(0); i < 8; i++ {
		accum += (n[i] - P64[i]) >> 64
	}

	return n, accum
}

func (n *bigNumber64) setBytes64(in []byte) *bigNumber64 {
	if len(in) != 56 {
		return nil
	}

	s := serialized{}
	for i, si := range in {
		s[len(s)-i-1] = si
	}

	d, _ := decafDeser64(s)

	for i, di := range d {
		n[i] = di
	}

	return n
}

// this will replace serialize
func decafSerialize64(dst []byte, n *bigNumber64) {
	n.decafCanon64()

	var bits uint
	var buf dword_t
	var k int

	for i := 0; i < 8; i++ {
		buf |= n[i] << bits

		for bits += 56; (bits >= 8 || i == 8-1) && k < 56; buf, bits, k = buf>>8, bits-8, k+1 {
			dst[k] = byte(buf) // why is msb set to 0 as default?
		}
	}
}
