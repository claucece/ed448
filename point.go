package ed448

import (
	"errors"
	"fmt"
)

var (
	bigNumOne           = mustDeserialize(serialized{1})
	bigNumTwo           = mustDeserialize(serialized{2})
	curveDSigned        = int64(-39081)
	twistedCurveDSigned = int64(-39082)
	sqrtDminus1         = mustDeserialize(serialized{
		0x46, 0x9f, 0x74, 0x36, 0x18, 0xe2, 0xd2, 0x79,
		0x01, 0x4f, 0x2b, 0xb4, 0x8d, 0x88, 0x38, 0xea,
		0xde, 0xab, 0x9a, 0x18, 0x5a, 0x06, 0x4c, 0xf1,
		0xa6, 0x5c, 0xe6, 0x51, 0x70, 0x97, 0x4d, 0x42,
		0x7b, 0x9f, 0xa4, 0x56, 0xf6, 0xc5, 0x28, 0x46,
		0xac, 0xdc, 0x4a, 0x73, 0x48, 0x87, 0x3b, 0x44,
		0x49, 0x7a, 0x5b, 0xb2, 0xc0, 0xc0, 0xfe, 0x12,
	})
)

func maskToBoolean(m uint32) bool {
	return m == 0xffffffff
}

// XXX: change name of all methods var

// NewPoint instantiates a new point in a suitable coordinate system.
// The x and y coordinates must be affine coordinates in little-endian
//XXX This should probably receive []byte{}

func NewPoint(x serialized, y serialized) (p *homogeneousProjective, e error) {
	xN, ok1 := deserialize(x)
	yN, ok2 := deserialize(y)

	p = newHomogeneousProjective(xN, yN)

	if !(ok1 && ok2) {
		e = errors.New("invalid coordinates")
	}

	return
}

type extensibleCoordinates struct {
	x, y, z, t, u *bigNumber
}

//Affina(x,y) => extensible(X, Y, Z, T, U)
func newExtensible(px, py *bigNumber) *extensibleCoordinates {
	x := px.copy()
	y := py.copy()
	z := &bigNumber{1}
	t := x.copy()
	u := y.copy()

	return &extensibleCoordinates{
		x: x,
		y: y,
		z: z,
		t: t,
		u: u,
	}
}

func (p *extensibleCoordinates) OnCurve() bool {
	x := p.x
	y := p.y
	z := p.z
	t := p.t
	u := p.u

	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)

	// Check invariant:
	// 0 = d*t^2*u^2 - x^2 - y^2 + z^2
	l2 = l2.square(y)
	l1 = l1.neg(l2)
	l0 = l0.square(z)
	l2 = l2.add(l0, l1)
	l3 = l3.square(u)
	l0 = l0.square(t)
	l1 = l1.mul(l0, l3)
	l0 = l0.mulWSignedCurveConstant(l1, curveDSigned)
	l1 = l1.add(l0, l2)
	l0 = l0.square(x)
	l2 = l2.neg(l0)
	l0 = l0.add(l2, l1)
	l5 := l0.zeroMask()

	// Check invariant:
	// 0 = -x*y + z*t*u
	l1 = l1.mul(t, u)
	l2 = l2.mul(z, l1)
	l0 = l0.mul(x, y)
	l1 = l1.neg(l0)
	l0 = l0.add(l1, l2)

	l4 := l0.zeroMask()

	ret := l4 & l5 & (^z.zeroMask())
	return maskToBoolean(ret)
}

func (p *extensibleCoordinates) equals(q *extensibleCoordinates) bool {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2 = l2.mul(q.z, p.x)
	l1 = l1.mul(p.z, q.x)
	l0 = l0.sub(l2, l1)
	l4 := l0.zeroMask()

	l2 = l2.mul(q.z, p.y)
	l1 = l1.mul(p.z, q.y)
	l0 = l0.sub(l2, l1)
	l3 := l0.zeroMask()

	return maskToBoolean(l4 & l3)
}

type twPNiels struct {
	n *twNiels
	z *bigNumber
}

func newTwistedPNiels(a, b, c, z [56]byte) *twPNiels {
	return &twPNiels{
		&twNiels{
			a: mustDeserialize(serialized(a)),
			b: mustDeserialize(serialized(b)),
			c: mustDeserialize(serialized(c)),
		},
		mustDeserialize(serialized(z)),
	}
}

func (p *twPNiels) String() string {
	return fmt.Sprintf("A: %s\nB: %s\nC: %s\nZ: %s\n", p.n.a, p.n.b, p.n.c, p.z)
}

func (p *twPNiels) equals(p2 *twPNiels) bool {
	ok := true

	ok = ok && p.n.equals(p2.n)
	ok = ok && p.z.equals(p2.z)

	return ok
}

type twNiels struct {
	a, b, c *bigNumber
}

func newNielsPoint(a, b, c [56]byte) *twNiels {
	return &twNiels{
		a: mustDeserialize(serialized(a)),
		b: mustDeserialize(serialized(b)),
		c: mustDeserialize(serialized(c)),
	}
}

func (p *twNiels) equals(p2 *twNiels) bool {
	ok := true

	ok = ok && p.a.equals(p2.a)
	ok = ok && p.b.equals(p2.b)
	ok = ok && p.c.equals(p2.c)

	return ok
}

func (p *twNiels) String() string {
	return fmt.Sprintf("A: %s\nB: %s\nC: %s\n", p.a, p.b, p.c)
}

func (p *twNiels) copy() *twNiels {
	return &twNiels{
		a: p.a.copy(),
		b: p.b.copy(),
		c: p.c.copy(),
	}
}

func (p *twNiels) conditionalNegate(neg word_t) {
	p.a.conditionalSwap(p.b, neg)
	p.c = p.c.conditionalNegate(neg)
}

func convertTwNielsToTwExtensible(dst *twExtensible, src *twNiels) {
	dst.y = dst.y.add(src.b, src.a)
	dst.x = dst.x.sub(src.b, src.a)
	dst.z = dst.z.setUi(1)
	dst.t = dst.x.copy()
	dst.u = dst.y.copy()
}

type twExtensible struct {
	x, y, z, t, u *bigNumber
}

func (p *twExtensible) copy(e *twExtensible) *twExtensible {
	p.x = e.x.copy()
	p.y = e.y.copy()
	p.z = e.z.copy()
	p.t = e.t.copy()
	p.u = e.u.copy()

	return p
}

func (p *twExtensible) addTwPNiels(a *twPNiels) *twExtensible {
	p.z.mulCopy(p.z, a.z)
	return p.addTwNiels(a.n)
}

func (p *twExtensible) subTwPNiels(a *twPNiels) {
	p.z.mulCopy(p.z, a.z)
	p.subTwNiels(a.n)
}

func convertTwExtensibleToTwPNiels(dst *twPNiels, src *twExtensible) {
	dst.n.a.sub(src.y, src.x)
	dst.n.b.add(src.x, src.y)
	dst.z.mul(src.u, src.t)
	dst.n.c.mulWSignedCurveConstant(dst.z, curveDSigned*2-2)
	dst.z.add(src.z, src.z)
}

func (p *twExtensible) twPNiels() *twPNiels {
	ret := &twPNiels{
		n: &twNiels{
			a: new(bigNumber),
			b: new(bigNumber),
			c: new(bigNumber),
		},
		z: new(bigNumber),
	}

	convertTwExtensibleToTwPNiels(ret, p)
	return ret
}

func convertTwPnielsToTwExtensible(dst *twExtensible, src *twPNiels) {
	dst.u.add(src.n.b, src.n.a)
	dst.t.sub(src.n.b, src.n.a)
	dst.x.mul(src.z, dst.t)
	dst.y.mul(src.z, dst.u)
	dst.z.square(src.z)
}

func (p *twExtensible) OnCurve() bool {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)

	// Check invariant:
	// 0 = -x*y + z*t*u
	l1 = l1.mul(p.t, p.u)
	l2 = l2.mul(p.z, l1)
	l0 = l0.mul(p.x, p.y)
	l1 = l1.neg(l0)
	l0 = l0.add(l1, l2)
	l5 := l0.zeroMask()

	// Check invariant:
	// 0 = d*t^2*u^2 + x^2 - y^2 + z^2 - t^2*u^2

	l2 = l2.square(p.y)
	l1 = l1.neg(l2)
	l0 = l0.square(p.x)
	l2 = l2.add(l0, l1)
	l3 = l3.square(p.u)
	l0 = l0.square(p.t)
	l1 = l1.mul(l0, l3)
	l3 = l3.mulWSignedCurveConstant(l1, curveDSigned)
	l0 = l0.add(l3, l2)
	l3 = l3.neg(l1)
	l2 = l2.add(l3, l0)
	l1 = l1.square(p.z)
	l0 = l0.add(l1, l2)
	l4 := l0.zeroMask()

	ret := l4 & l5 & (^p.z.zeroMask())
	return maskToBoolean(ret)
}

func (p *twExtensible) setIdentity() {
	p.x.setUi(0)
	p.y.setUi(1)
	p.z.setUi(1)
	p.t.setUi(0)
	p.u.setUi(0)
}

func (p *twExtensible) String() string {
	ret := fmt.Sprintf("X: %s\n", p.x)
	ret += fmt.Sprintf("Y: %s\n", p.y)
	ret += fmt.Sprintf("Z: %s\n", p.z)
	ret += fmt.Sprintf("T: %s\n", p.t)
	ret += fmt.Sprintf("U: %s\n", p.u)

	return ret
}

func (p *twExtensible) equals(p2 *twExtensible) bool {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2 = l2.mul(p2.z, p.x)
	l1 = l1.mul(p.z, p2.x)
	l0 = l0.sub(l2, l1)

	l4 := l0.zeroMask()

	l2 = l2.mul(p2.z, p.y)
	l1 = l1.mul(p.z, p2.y)
	l0 = l0.sub(l2, l1)

	l3 := l0.zeroMask()

	return (l4 & l3) == 0xffffffff
}

func (p *twExtensible) double() *twExtensible {
	x := p.x
	y := p.y
	z := p.z
	t := p.t
	u := p.u

	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2.square(x)
	l0.square(y)
	u = u.addRaw(l2, l0)
	t = t.addRaw(y, x)
	l1.square(t)
	t = t.subRaw(l1, u)
	t.bias(3)
	t.weakReduce()
	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	l1 = l1.sub(l0, l2)
	x.square(z)
	x.bias(1)
	z = z.addRaw(x, x)
	l0 = l0.subRaw(z, l1)
	l0.weakReduce()
	z.mul(l1, l0)
	x.mul(l0, t)
	y.mul(l1, u)

	return p
}

func (p *twExtensible) addTwNiels(p2 *twNiels) *twExtensible {
	x := p.x
	y := p.y
	z := p.z
	t := p.t
	u := p.u

	l0 := new(bigNumber)
	l1 := new(bigNumber)

	l1 = l1.sub(y, x)
	l0.mul(p2.a, l1)
	l1 = l1.addRaw(x, y)
	y.mul(p2.b, l1)
	l1.mul(u, t)
	x.mul(p2.c, l1)

	u = u.addRaw(l0, y)
	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	t = t.sub(y, l0)

	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	y = y.sub(z, x)
	l0 = l0.addRaw(x, z)

	z.mul(l0, y)
	x.mul(y, t)
	y.mul(l0, u)

	return p
}

func (p *twExtensible) subTwNiels(e *twNiels) {
	L1 := new(bigNumber).subxRaw(p.y, p.x)
	L0 := new(bigNumber).mul(e.b, L1)
	L1.addRaw(p.x, p.y)
	p.y.mul(e.a, L1)
	L1.mul(p.u, p.t)
	p.x.mul(e.c, L1)
	p.u.addRaw(L0, p.y)
	p.t.subxRaw(p.y, L0)
	p.y.addRaw(p.x, p.z)
	L0.subxRaw(p.z, p.x)
	p.z.mul(L0, p.y)
	p.x.mul(p.y, p.t)
	p.y.mul(L0, p.u)
}

func (p *twExtensible) untwistAndDoubleAndSerialize() *bigNumber {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)
	b := new(bigNumber)

	l3.mul(p.y, p.x)
	b.add(p.y, p.x)
	l1.square(b)
	l2.add(l3, l3)
	b.sub(l1, l2)
	l2.square(p.z)
	l1.square(l2)
	b.add(b, b)
	l2.mulWSignedCurveConstant(b, curveDSigned-1)
	b.mulWSignedCurveConstant(l2, curveDSigned-1)
	l0.mul(l2, l1)
	l2.mul(b, l0)
	l0.isr(l2)
	l1.mul(b, l0)

	//XXX This is included in the original code, but it seems not to be used
	//b = b.square(l0)
	//l0 = l0.mul(l2, b)

	return b.mul(l1, l3)
}

type pointT struct {
	x, y, z, t *bigNumber
}

func hibit(x *bigNumber) word_t {
	y := &bigNumber{}
	y.add(x, x)
	y.decafCanon()
	return word_t(-(y[0] & 1))
}

// this is replacing untwistAndSerialize method and serialize
// XXX: check and compare with strike' functions and fast_decaf
// xxx: add sub is the same
// xxx : mul and mulW are not the same --> check
func (p *pointT) desisogenize() *bigNumber {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	a.decafMulW(p.y, 1-(-39081))
	c.decafMul(a, p.t) // maybe b
	a.decafMul(p.x, p.z)
	d.decafSub(c, a) // s := |(u . (r . (aZ . X-d . Y . T) + Y ) /a|
	a.decafAdd(p.z, p.y)
	b.decafSub(p.z, p.y)
	c.decafMul(b, a)
	b.decafMulW(c, (-(-39081)))
	a.isr(b)                         // r := 1/sqrt((a-d) . (Z+X) . (Z-Y))
	b.decafMulW(a, (-(-36081)))      // u := (a - d) . r
	c.decafMul(b, a)                 // u . r
	a.decafMul(c, d)                 // (ur) . (aZT-dYT)
	d.decafAdd(b, b)                 // 2u = -2au since a = -1
	c.decafMul(d, p.z)               // 2u . Z
	b.conditionalNegate(^(hibit(c))) // u := -u if negative
	c.decafMul(b, p.y)               // final y?
	a.decafAdd(a, c)
	a.conditionalNegate(hibit(a)) // a?

	return a

}

//HP(X : Y : Z) = Affine(X/Z, Y/Z), Z ≠ 0
//XXX This can be replaced by extensible for simplicity if we neither use ADD
//on the basePoint in test and benchmark (it is not used elsewhere)
type homogeneousProjective struct {
	x, y, z *bigNumber
}

//Affine to Homogeneous Projective
func newHomogeneousProjective(x *bigNumber, y *bigNumber) *homogeneousProjective {
	return &homogeneousProjective{
		x: x.copy(),
		y: y.copy(),
		z: bigNumOne.copy(),
	}
}

func (p *homogeneousProjective) String() string {
	return fmt.Sprintf("X: %s\nY: %s\nZ: %s\n", p.x, p.y, p.z)
}

func (p *homogeneousProjective) OnCurve() bool {
	// (x² + y²)z² - z^4 - dx²y² = 0
	x := p.x
	y := p.y
	z := p.z

	x2 := new(bigNumber).mul(x, x)
	y2 := new(bigNumber).mul(y, y)
	z2 := new(bigNumber).mul(z, z)
	z4 := new(bigNumber).mul(z2, z2)

	x2y2 := new(bigNumber).mul(x2, y2)
	dx2y2 := x2y2.mulWSignedCurveConstant(x2y2, curveDSigned)
	dx2y2.weakReduce()

	r := new(bigNumber).add(x2, y2)
	r.mulCopy(r, z2)
	r.sub(r, z4)
	r.sub(r, dx2y2)

	r.strongReduce()
	return r.zero()
}

func rev(in []byte) []byte {
	r := make([]byte, len(in), len(in))

	for i, ni := range in {
		r[len(in)-i-1] = ni
	}

	return r
}

// See Hisil, formula 5.1
//XXX Used only for testing
func (p *homogeneousProjective) double() *homogeneousProjective {
	x1 := p.x
	y1 := p.y
	z1 := p.z

	b := new(bigNumber).square(new(bigNumber).add(x1, y1))
	c := new(bigNumber).square(x1)
	d := new(bigNumber).square(y1)
	e := new(bigNumber).add(c, d)
	h := new(bigNumber).square(z1)
	//j := h.mulW(h, 2) // This is slower than adding
	j := h.add(h, h)
	j.sub(e, j)

	xx := b.sub(b, e)
	xx.mulCopy(xx, j)
	yy := c.sub(c, d)
	yy.mulCopy(yy, e)
	zz := e.mulCopy(e, j)

	//XXX PERF Should it change the same instance instead?
	return &homogeneousProjective{
		xx, yy, zz,
	}
}

// See Hisil, formula 5.3
func (p *homogeneousProjective) add(p2 *homogeneousProjective) *homogeneousProjective {
	//A ← Z1*Z2,
	//B ← A^2,
	//C ← X1*X2,
	//D ← Y1*Y2,
	//E ← dC*D,
	//F ← B−E,
	//G ← B+E,
	//X3 ← A*F*((X1+Y1)*(X2+Y2)−C−D),
	//Y3 ← A*G*(D−aC),
	//Z3 ← F*G.

	x1 := p.x
	y1 := p.y
	z1 := p.z

	x2 := p2.x
	y2 := p2.y
	z2 := p2.z

	a := new(bigNumber).mul(z1, z2)
	b := new(bigNumber).square(a)
	c := new(bigNumber).mul(x1, x2)
	d := new(bigNumber).mul(y1, y2)

	e := new(bigNumber).mulWSignedCurveConstant(c, curveDSigned)
	e.mulCopy(e, d)
	f := new(bigNumber).sub(b, e)
	g := new(bigNumber).add(b, e)

	//Just reusing e and b (unused) memory
	x3 := e.mulCopy(b.add(x1, y1), e.add(x2, y2))
	x3.sub(x3, c).sub(x3, d)
	x3.mulCopy(x3, a).mulCopy(x3, f)

	y3 := d.sub(d, c)
	y3 = y3.mulCopy(y3, a).mulCopy(y3, g)

	z3 := f.mulCopy(f, g)

	return &homogeneousProjective{
		x3, y3, z3,
	}
}

//XXX Move: bigNumber should not know about points
func (sz *bigNumber) deserializeAndTwistApprox() (*twExtensible, bool) {
	a := &twExtensible{
		x: new(bigNumber),
		y: new(bigNumber),
		z: new(bigNumber),
		u: new(bigNumber),
		t: new(bigNumber),
	}

	var L0, L1 *bigNumber
	L0 = new(bigNumber)
	L1 = new(bigNumber)
	a.z.square(sz)
	a.y = a.z.copy()
	a.y.addW(1)
	L0.square(a.y)
	a.x.mulWSignedCurveConstant(L0, curveDSigned-1)
	a.y.add(a.z, a.z)
	a.u.add(a.y, a.y)
	a.y.add(a.u, a.x)
	a.x.square(a.z)
	a.u.neg(a.x)
	a.u.addW(1)
	a.x.mul(sqrtDminus1, a.u)
	L0.mul(a.x, a.y)
	a.t.mul(L0, a.y)
	a.u.mul(a.x, a.t)
	a.t.mul(a.u, L0)
	a.y.mul(a.x, a.t)
	L0.isr(a.y)
	a.y.mul(a.u, L0)
	L1.square(L0)
	a.u.mul(a.t, L1)
	a.t.mul(a.x, a.u)
	a.x.add(sz, sz)
	L0.mul(a.u, a.x)
	a.x = a.z.copy()
	L1.neg(a.x)
	L1.addW(1)
	a.x.mul(L1, L0)
	L0.mul(a.u, a.y)
	a.z.addW(1)
	a.y.mul(a.z, L0)
	a.t.subW(1)

	// XXX maybe related with constant time
	ret := a.t.zero()

	a.z.setUi(1)
	a.t = a.x.copy()
	a.u = a.y.copy()

	return a, !ret
}
