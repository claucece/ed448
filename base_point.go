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
