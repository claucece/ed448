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

//This function runs Elligator2 on the decaf Jacobi quartic model.  It then
// uses the isogeny to put the result in twisted Edwards form.  As a result,
// it is safe (cannot produce points of order 4), and would be compatible with
// hypothetical other implementations of Decaf using a Montgomery or untwisted
// Edwards model.
func decafNonuniformMapToCurve(p *pointT, ser serialized) {
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
	urr, a, b, c, dee, e, ur2D, udr21 := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	r, _ := decafDeser(ser)
	r.decafCanon()
	a.decafSqr(r)                         //r^2
	urr.decafMulW(a, QuadraticNonresidue) // urr = u*r^2
	dee.decafMulW(&bigNumber{0x01}, D)    // dee = 1*D
	a.decafAdd(urr, &bigNumber{0x01})     // u*r^2 + 1
	ur2D.decafSub(dee, urr)               // ur2_d = -(u*r^2-d)
	c.decafMul(a, ur2D)                   // (r^2 * -(u*r^2-d))
	b.decafMulW(urr, -(D))                // (u*r^2 -d)
	udr21.decafAdd(b, &bigNumber{0x01})   // udr21 = -(udr^2-1)
	a.decafMul(c, udr21)                  // (r^2 * -(u*r^2-d)) * -(udr^2-1)
	c.decafMulW(a, D+1)                   // c = (u*r^2 + 1) * (d - u*r^2) * (1 - u*d*r^2) * (d+1)
	b.decafIsqrt(c)                       // FIELD: if 5 mod 8, multiply result by u (aka urr)
	a.decafSqr(b)                         // (u*r^2 -d)^2
	e.decafMul(a, c)                      // (u*r^2 -d)^2 * (u*r^2 + 1) * (d - u*r^2) * (1 - u*d*r^2) * (d+1)

	mask := decafEq(e, &bigNumber{0x01}) // mask for trailling zeros
	a.decafMul(b, r)                     // (u*r^2 -d) * r
	b.decafCondSel(a, b, mask)           // mask? a : b
	b.decafCondNegate(hibit(b))          //-b
	a.decafMulW(b, D+1)                  //-b * D+1

	/* Here: a = sqrt( (d+1) / (ur^2?) * (u*r^2 + 1) * (d - u*r^2) * (1 - u*d*r^2)) */

	ur2D.decafCondSwap(udr21, ^(mask)) //
	e.decafMul(ur2D, a)                // (-(u*r^2-d)) * (u*r^2 -d)^2
	b.decafMul(udr21, a)               // -(udr^2-1) * (u*r^2 -d)^2
	c.decafSqr(b)                      // square above

	/* Here:
	 * ed_x = 2e/(1-e^2)
	 * c =  * (ed_y-1)/(ed_y+1)
	 *
	 * Special cases:
	 *   e^2 = 1: impossible for cofactor-4 curves (would isogenize to order-4 point)
	 *   e = 0 <-> also c = 0: maps to (0,1), which is fine.
	 */

	a.decafSqr(e)
	a.decafSub(&bigNumber{0x01}, a)
	e.decafAdd(e, e)
	b.decafAdd(dee, c)
	c.decafSub(dee, c)

	p.x.decafMul(e, c)
	p.z.decafMul(a, c)
	p.y.decafMul(b, a)
	p.t.decafMul(b, e)

}
