package ed448

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
