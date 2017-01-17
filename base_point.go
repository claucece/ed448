package ed448

func (p *pointT) decafPointValidate() word_t {
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

// play with this
//func (p *CompletedGroupElement) ToExtended(r *ExtendedGroupElement) {
//	FeMul(&r.X, &p.X, &p.T)
//	FeMul(&r.Y, &p.Y, &p.Z)
//	FeMul(&r.Z, &p.Z, &p.T)
//	FeMul(&r.T, &p.X, &p.Y)
//}
