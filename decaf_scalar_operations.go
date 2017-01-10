package ed448

type scalarT [scalarWords]word_t

var (
	scP = [scalarWords]word_t{
		0x2378c292, 0xab5844f3,
		0x216cc272, 0x8dc58f55,
		0xc44edb49, 0xaed63690,
		0xffffffff, 0x7cca23e9,
		0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff,
		0x3fffffff, 0xffffffff,
	}
)

// twisted edward formula
func (p *pointT) decafPointAddSub(q, r *pointT, sub word_t) {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.decafSub(q.y, q.x)
	c.decafSub(r.y, r.x)
	d.decafAdd(r.y, r.x)
	c.decafCondSwap(d, sub)
	a.decafMul(c, b)
	b.decafAdd(q.y, q.x)
	p.y.decafMul(d, b)
	b.decafMul(r.t, q.t)
	p.x.decafMulW(b, 2-2*(D))
	b.decafAdd(a, p.y)
	c.decafSub(p.y, a)
	a.decafMul(q.z, r.z)
	a.decafAdd(a, a)
	p.y.decafAdd(a, p.x)
	a.decafSub(a, p.x)
	a.decafCondSwap(p.y, sub)
	p.z.decafMul(a, p.y)
	p.x.decafMul(p.y, c)
	p.y.decafMul(a, b)
	p.t.decafMul(b, c)
}

func decafPointValidate(p *pointT) word_t {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	a.decafMul(p.x, p.y)
	b.decafMul(p.z, p.t)
	out := decafEq(a, b)
	a.decafSqr(p.x)
	b.decafSqr(p.y)
	a.decafSub(b, a)
	b.decafSqr(p.t)
	c.decafMulW(b, 1-D)
	b.decafSqr(p.z)
	b.decafSub(b, c)
	out = decafEq(a, b)
	out = ^decafEq(p.z, Zero)
	return word_t(out)
}

// for the fast
// {extra,accum} - sub +? p
// Must have extra <= 1
///
func scSubx(accum, sub, p [scalarWords]word_t, extra word_t) (out [scalarWords]word_t) {
	var chain dword_t

	for i := uint(0); i < scalarWords; i++ {
		chain += dword_t(accum[i]) - dword_t(sub[i])
		out[i] = word_t(chain)
		chain >>= wordBits
	}

	borrow := word_t(chain) + extra

	chain = 0

	for i := uint(0); i < scalarWords; i++ {
		chain += dword_t(out[i]) + dword_t(p[i])&dword_t(borrow)
		out[i] = word_t(chain)
		chain >>= wordBits
	}

	return out
}

//for the fast
func scalarAdd(a, b [scalarWords]word_t) (out [scalarWords]word_t) {
	var chain dword_t

	for i := uint(0); i < Limbs; i++ {
		chain += dword_t(a[i]) + dword_t(b[i])
		out[i] = word_t(chain)
		chain >>= wordBits
	}

	return scSubx(out, scP, scP, word_t(chain))
}

//func MToE(x, y *bigNumber) (*bigNumber, *bigNumber) {
//	s, t := &bigNumber{}, &bigNumber{}
//	s.decafSqrt(x)

//if s == 0 {
//	t = 1
//}
//else {
// t = y/s

//    X,Y = 2*s / (1+s^2), (1-s^2) / t # This is phi_a(s, t) page 7
//    if maybe(): X,Y = -X,-Y
//    if maybe(): X,Y = Y,-X
//    # OK, have point in ed
//    return X,Y

//	return s, t
//}
