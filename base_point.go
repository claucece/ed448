package ed448

const (
	lmask = (dword_t(1) << 56) - 1
)

type bigNumber64 [8]dword_t

type base struct {
	x, y, z, t *bigNumber64
}

func (n *bigNumber64) copy64() *bigNumber64 {
	c := &bigNumber64{}
	copy(c[:], n[:])
	return c
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
