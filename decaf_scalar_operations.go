package ed448

type scalarT [scalarWords]word_t

var (
	scP = [scalarWords]word_t{
		0xab5844f3,
		0x2378c292,
		0x8dc58f55,
		0x216cc272,
		0xaed63690,
		0xc44edb49,
		0x7cca23e9,
		0xffffffff,
		0xffffffff,
		0xffffffff,
		0xffffffff,
		0xffffffff,
		0xffffffff,
		0x3fffffff,
	}
)

// twisted edward formula
// from the normal decaf
// XXX: decide which one is going to be used
func (p *pointT) decafPointAddSub(q, r *pointT, sub dword_t) {
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

func (p *pointT) decafPointAdd(a, b *pointT) {
	p.decafPointAddSub(a, b, dword_t(0))
}

func (p *pointT) decafPointSub(a, b *pointT) {
	p.decafPointAddSub(a, b, dword_t(0xffffffff))
}

// from now on this is decaf_fast
// from decaf_fast
// {extra,accum} - sub + p
// Must have extra <= 1
// XXX: check if this function is doing exactly the same as STRIKE one
func scSubx(accum, sub, p [scalarWords]word_t, extra word_t) (out [scalarWords]word_t) {
	var chain int64

	for i := uint(0); i < scalarWords; i++ {
		chain += int64(accum[i]) - int64(sub[i])
		out[i] = word_t(chain)
		chain >>= wordBits
	}

	borrow := chain + int64(extra) // 0 or -1
	chain = 0

	for i := uint(0); i < scalarWords; i++ {
		chain += int64(out[i]) + (int64(p[i]) & borrow)
		out[i] = word_t(chain)
		chain >>= wordBits
	}
	return out
}

func scalarAdd(a, b [scalarWords]word_t) (out [scalarWords]word_t) {
	var chain dword_t

	for i := uint(0); i < scalarWords; i++ {
		chain += dword_t(a[i]) + dword_t(b[i])

		out[i] = word_t(chain)
		chain >>= wordBits
	}

	return scSubx(out, scP, scP, word_t(chain))
}

func scalarSub(a, b [scalarWords]word_t) (out [scalarWords]word_t) {
	return scSubx(a, b, scP, 0)
}

func scalarCopy(a [scalarWords]word_t) (out [scalarWords]word_t) {
	copy(out[:], a[:])
	return out
}

func scHalve(a, b [scalarWords]word_t) (out [scalarWords]word_t) {
	mask := -(a[0] & 1)
	var chain dword_t
	var i uint

	for i = 0; i < scalarWords; i++ {
		chain += dword_t(a[i]) + dword_t(b[i]&mask)
		out[i] = word_t(chain)
		chain >>= wordBits
	}

	for i = 0; i < scalarWords-1; i++ {
		out[i] = out[i]>>1 | out[i+1]<<(wordBits-1)
	}

	out[i] = out[i]>>1 | word_t(chain<<(wordBits-1))

	return
}

//In Progress
func scalarAdjustment() [scalarWords]word_t {
	var smadj [scalarWords]word_t
	one := [scalarWords]word_t{0x01}
	smadj = scalarCopy(one)

	for i := uint(0); i < uint(8)*uint(4)*uint(14); i++ { //works until 8*8*13
		smadj = scalarAdd(smadj, smadj)
	}
	smadj = scalarSub(smadj, one)
	return smadj
}

func (p *pointT) addNielsToProjective(p2 *twNiels, beforeDouble bool) {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.sub(p.y, p.x)
	a.mul(p2.a, b)
	b.addRaw(p.x, p.y)
	p.y.mul(p2.b, b)
	p.x.mul(p2.c, p.t)
	c.addRaw(a, p.y)
	b.sub(p.y, a)
	p.y.sub(p.z, p.x)
	a.addRaw(p.x, p.z)
	p.z.mul(a, p.y)
	p.x.mul(p.y, b)
	p.y.mul(a, c)
	if !beforeDouble {
		p.t.mul(b, c)
	}
}

func convertNielsToPt(dst *pointT, src *twNiels) {
	dst.y.add(src.b, src.a)
	dst.x.sub(src.b, src.a)
	dst.t.mul(dst.y, dst.x)
	dst.z.copyFrom(One)
}

// Based on Hisil's formula 5.1.3: Doubling in E^e
func (p *pointT) pointDoubleInternal(q *pointT, beforeDouble bool) {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	c.square(q.x)
	a.square(q.y)
	d.addRaw(c, a)
	p.t.addRaw(q.y, q.x)
	b.square(p.t)
	exponentBias := uint32(3)
	b.decafSubRawWithX(b, d, exponentBias)
	p.t.sub(a, c)
	p.x.square(q.z)
	p.z.addRaw(p.x, p.x)
	exponentBias = uint32(4)
	a.decafSubRawWithX(p.z, p.t, exponentBias)
	p.x.mul(a, b)
	p.z.mul(p.t, a)
	p.y.mul(p.t, d)
	if !beforeDouble {
		p.t.mul(b, d)
	}
}

func (c *curveT) precomputedScalarMul(scalar [scalarWords]word_t) *pointT {
	out := &pointT{
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
	}

	n := uint(5)
	t := uint(5)
	s := uint(18)

	var scalar1 [scalarWords]word_t
	scalar1 = scalarAdd(scalar, precomputedBaseTable.scalarAdjustment)

	scalar1 = scHalve(scalar1, scP)

	var ni *twNiels

	for i := int(s - 1); i >= 0; i-- {
		if i != int(s-1) {
			out.pointDoubleInternal(out, false)
		}

		for j := uint(0); j < n; j++ {
			var tab word_t
			for k := uint(0); k < t; k++ {
				bit := uint(i) + s*(k+j*t)
				if bit < 446 { // change 446 to constant
					tab |= (scalar1[bit/uint(32)] >> (bit % uint(32)) & 1) << k
					// change uint(32) to constant
				}
			}

			invert := (int32(tab) >> (t - 1)) - 1
			tab ^= word_t(invert)
			tab &= (1 << (t - 1)) - 1

			ni = precomputedBaseTable.decafLookup(j, t, uint(tab))

			ni.conditionalNegate(word_t(invert))

			if i != int(s-1) || j != 0 {
				out.addNielsToProjective(ni, j == n-1 && i != 0)
			} else {
				convertNielsToPt(out, ni)
			}
		}
	}

	return out
}
