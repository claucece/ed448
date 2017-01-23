package ed448

// q is a odd prime power
// u : a non square in F_q
// A, B := elements of Fq such that AB(A^2 - 4B) != 0
// R := {r in F_q : 1 + ur^2 != 0, A^2ur^2 != B(1 + ur^2)^2

// v := -A/(1/ur^2)
// ep := x_(v^3+Av^2+Bv)
// x:= epv - (1 - ep)A/2
// y := -ep * sqrt(x^3_Ax^2_Bx)

// x_ :=  : F_q -> F_q by x_(a)^((q-1)/2)
//u = (1 - t)/(1 + t)
//v = u^5 + (r^2 - 2)u^3 + u
//X = x_(v)u
//Y = (x_(v)v^((q+1)/4) x_(v)*x_(u^2+1/c^2)
//x = (c - 1)sX(1 + X)/Y
//y = (rX -(1 + X)^2) = (rX + (1 + X)^2)

const (
	QuadraticNonresidue = -1
)

/** Inverse square root using addition chain. */
func decafIsqrtChk(y, x *bigNumber, zero dword_t) dword_t {
	tmp0, tmp1 := &bigNumber{}, &bigNumber{}
	y.decafIsqrt(x)
	tmp0.decafSqr(y)
	tmp1.decafMul(tmp0, x)
	return decafEq(tmp1, &bigNumber{0x01}) | (zero & decafEq(tmp1, &bigNumber{0x00}))
}

// 2-torque a point
func point2Torque(p, q *pointT) {
	p.x.decafSub(&bigNumber{0x00}, q.x)
	p.y.decafSub(&bigNumber{0x00}, q.y)
	p.z.copyFrom(q.z)
	p.t.copyFrom(q.t)

}

//This function runs Elligator2 on the decaf Jacobi quartic model.  It then
// uses the isogeny to put the result in twisted Edwards form.  As a result,
// it is safe (cannot produce points of order 4), and would be compatible with
// hypothetical other implementations of Decaf using a Montgomery or untwisted
// Edwards model.
// gives out the data hashed to the curve
func decafNonuniformMapToCurve(ser serialized) (*pointT, dword_t) {
	/*
	   XXD = (u*r^2 + 1) * (d - u*r^2) * (1 - u*d*r^2) / (d+1) // c=u*r^2 && s = r
	   factor(XX / (1/XXD))
	   (u*r^2 - d)^2
	   factor((ey-1)/(ey+1)/(1/d * 1/XXD))
	   (u*d*r^2 - 1)^2
	   factor(XX2 / (u*r^2/XXD))
	   (u*d*r^2 - 1)^2
	   factor((ey2-1)/(ey2+1)/(1/d * u*r^2/XXD))
	   (u*r^2 - d)^2
	*/
	r, a, b, c, dee, d2, n, rn, e := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	p := &pointT{
		x: new(bigNumber),
		y: new(bigNumber),
		z: new(bigNumber),
		t: new(bigNumber),
	}

	// probable nonresidue
	r0, overT := (decafDeser(ser))
	over := ^overT
	sgnR0 := hibit(r0)
	r0.decafCanon()
	a.decafSqr(r0) //r^2
	r.decafSub(&bigNumber{0x00}, a)
	//r.decafMulW(a, QuadraticNonresidue) // urr = u*r^2
	dee.decafMulW(&bigNumber{0x01}, D) // dee = 1*D
	c.decafMul(r, d2)

	//a.decafAdd(urr, &bigNumber{0x01})             // u*r^2 + 1
	//ur2D.decafSub(dee, urr)                       // ur2_d = -(u*r^2-d)
	//c.decafMul(a, ur2D)                           // (r^2 * -(u*r^2-d))
	//b.decafMulW(urr, -(D))                        // (u*r^2 -d)
	//udr21.decafAdd(b, &bigNumber{0x01})           // udr21 = -(udr^2-1)
	//a.decafMul(c, udr21)                          // (r^2 * -(u*r^2-d)) * -(udr^2-1)
	//c.decafMulW(a, D+1)                           // c = (u*r^2 + 1) * (d - u*r^2) * (1 - u*d*r^2) * (d+1)
	//b.decafIsqrt(c)                               // FIELD: if 5 mod 8, multiply result by u (aka urr)
	//a.decafSqr(b)                                 // (u*r^2 -d)^2
	//e.decafMul(a, c)                              // (u*r^2 -d)^2 * (u*r^2 + 1) * (d - u*r^2) * (1 - u*d*r^2) * (d+1)

	/* Compute D2 := (dr+a-d)(dr-ar-d) with a=1 */ // from Decaf paper
	a.decafSub(c, dee)                             // D - D
	a.decafAdd(a, &bigNumber{0x01})                // D + 1
	specialIdentity := decafEq(a, &bigNumber{0x00})
	b.decafSub(c, r)
	b.decafSub(b, dee)
	d2.decafMul(a, b)

	/* compute N := (r+1)(a-2d) */
	a.decafAdd(r, &bigNumber{0x01})
	n.decafMulW(a, 1-2*D)

	/* e = +-1/sqrt(+-ND) */
	rn.decafMul(r, n)
	a.decafMul(rn, d2)
	square := decafIsqrtChk(e, a, dword_t(0))
	isZero := decafEq(r, &bigNumber{0x00})
	square |= isZero
	square |= specialIdentity

	/* b <- t/s */
	c.decafCondSel(r0, r, square) /* r? = sqr ? r : 1 */

	/* In two steps to avoid overflow on 32-bit arch */
	a.decafMulW(c, 1-2*D)
	b.decafMulW(a, 1-2*D)
	c.decafSub(r, &bigNumber{0x01})
	a.decafMul(b, c) /* = r? * (r-1) * (a-2d)^2 with a=1 */
	b.decafMul(a, e)
	b.decafCondNegate(^square)
	c.decafCondSel(r0, &bigNumber{0x01}, square)
	a.decafMul(e, c)
	c.decafMul(a, d2) // 1/s except for sign.
	b.decafSub(b, c)

	/* a <- s = e * N * (sqr ? r : r0)
	 * e^2 r N D = 1
	 * 1/s =  1/(e * N * (sqr ? r : r0)) = e * D * (sqr ? 1 : r0)
	 */
	a.decafMul(n, r0)
	rn.decafCondSel(a, rn, square)
	a.decafMul(rn, e)
	c.decafMul(a, b)

	/* Normalize/negate */
	negS := hibit(a) ^ (^square) //not?
	a.decafCondNegate(negS)      /* ends up negative if ~square */
	sgnOverS := hibit(b) ^ negS
	sgnOverS &= ^decafEq(n, &bigNumber{0x00}) // not?
	sgnOverS |= decafEq(d2, &bigNumber{0x00})

	/* b <- t */
	tmp := decafEq(c, &bigNumber{0x00})
	b.decafCondSel(c, &bigNumber{0x01}, tmp) /* 0,0 -> 1,0 */

	/* isogenize */
	c.decafSqr(a)    /* s^2 */
	a.decafAdd(a, a) /* 2s */
	e.decafAdd(c, &bigNumber{0x01})
	p.t.decafMul(a, e) /* 2s(1+s^2) */
	p.x.decafMul(a, b) /* 2st */
	a.decafSub(&bigNumber{0x01}, c)
	p.y.decafMul(e, a) /* (1+s^2)(1-s^2) */
	p.z.decafMul(a, b) /* (1-s^2)t */

	succ := (^square & 1) | (sgnOverS & 2) | (sgnR0 & 4) | (over & 8)
	return p, succ
	//mask := decafEq(e, &bigNumber{0x01}) // mask for trailling zeros
	//	a.decafMul(b, r)                     // (u*r^2 -d) * r
	//	b.decafCondSel(a, b, mask)           // mask? a : b
	//	b.decafCondNegate(hibit(b))          //-b
	//	a.decafMulW(b, D+1)                  //-b * D+1
	//
	//	/* Here: a = sqrt( (d+1) / (ur^2?) * (u*r^2 + 1) * (d - u*r^2) * (1 - u*d*r^2)) */
	//
	//	ur2D.decafCondSwap(udr21, ^(mask)) //
	//	e.decafMul(ur2D, a)                // (-(u*r^2-d)) * (u*r^2 -d)^2
	//	b.decafMul(udr21, a)               // -(udr^2-1) * (u*r^2 -d)^2
	//	c.decafSqr(b)                      // square above
	//
	//	/* Here:
	//	 * ed_x = 2e/(1-e^2)
	//	 * c =  * (ed_y-1)/(ed_y+1)
	//	 *
	//	 * Special cases:
	//	 *   e^2 = 1: impossible for cofactor-4 curves (would isogenize to order-4 point)
	//	 *   e = 0 <-> also c = 0: maps to (0,1), which is fine.
	//	 */
	//
	//	a.decafSqr(e)
	//	a.decafSub(&bigNumber{0x01}, a)
	//	e.decafAdd(e, e)
	//	b.decafAdd(dee, c)
	//	c.decafSub(dee, c)
	//
	//	p.x.decafMul(e, c)
	//	p.z.decafMul(a, c)
	//	p.y.decafMul(b, a)
	//	p.t.decafMul(b, e)
	//
}
