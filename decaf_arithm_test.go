package ed448

import (
	_ "encoding/hex"

	. "gopkg.in/check.v1"
)

type DecafSuite struct{}

var _ = Suite(&DecafSuite{})

func (s *DecafSuite) Test_DecafCopy(c *C) {

	n := &bigNumber{
		0xaf1b9c5, 0xe417cd7, 0x839a472, 0x43bcebc,
		0x7dcf7e0, 0x726193c, 0x304a5bb, 0xf04e22e,
		0x526f560, 0x44604e2, 0x54f3a45, 0x979c291,
	}

	a := n.copy()

	c.Assert(a, DeepEquals, n)

}

func (s *DecafSuite) Test_DecafMul(c *C) {
	n := &bigNumber{}

	x, _ := deserialize(serialized{
		0xf5, 0x81, 0x74, 0xd5, 0x7a, 0x33, 0x72,
		0x36, 0x3c, 0x0d, 0x9f, 0xcf, 0xaa, 0x3d,
		0xc1, 0x8b, 0x1e, 0xff, 0x7e, 0x89, 0xbf,
		0x76, 0x78, 0x63, 0x65, 0x80, 0xd1, 0x7d,
		0xd8, 0x4a, 0x87, 0x3b, 0x14, 0xb9, 0xc0,
		0xe1, 0x68, 0x0b, 0xbd, 0xc8, 0x76, 0x47,
		0xf3, 0xc3, 0x82, 0x90, 0x2d, 0x2f, 0x58,
		0xd2, 0x75, 0x4b, 0x39, 0xbc, 0xa8, 0x74,
	})

	y, _ := deserialize(serialized{
		0x74, 0xa8, 0xbc, 0x39, 0x4b, 0x75, 0xd2,
		0x58, 0x2f, 0x2d, 0x90, 0x82, 0xc3, 0xf3,
		0x47, 0x76, 0xc8, 0xbd, 0x0b, 0x68, 0xe1,
		0xc0, 0xb9, 0x14, 0x3b, 0x87, 0x4a, 0xd8,
		0x7d, 0xd1, 0x80, 0x65, 0x63, 0x78, 0x76,
		0xbf, 0x89, 0x7e, 0xff, 0x1e, 0x8b, 0xc1,
		0x3d, 0xaa, 0xcf, 0x9f, 0x0d, 0x3c, 0x36,
		0x72, 0x33, 0x7a, 0xd5, 0x74, 0x81, 0xf5,
	})

	z, _ := deserialize(serialized{
		0x11, 0x95, 0x9c, 0x2e, 0x91, 0x78, 0x6f,
		0xec, 0xff, 0x37, 0xe5, 0x8e, 0x2b, 0x50,
		0x9e, 0xf8, 0xfb, 0x41, 0x08, 0xc4, 0xa7,
		0x02, 0x1c, 0xbf, 0x5a, 0x9f, 0x18, 0xa7,
		0xec, 0x32, 0x65, 0x7e, 0xed, 0xdc, 0x81,
		0x81, 0x80, 0xa8, 0x4c, 0xdd, 0x95, 0x14,
		0xe6, 0x67, 0x26, 0xd3, 0xa1, 0x22, 0xdb,
		0xb3, 0x9f, 0x17, 0x7a, 0x85, 0x16, 0x6c,
	})

	n.decafMul(x, y)

	c.Assert(n, DeepEquals, z)
}

func (s *DecafSuite) Test_DecafSqr(c *C) {
	n := &bigNumber{}

	x, _ := deserialize(serialized{0x08, 0xfd})

	m, _ := deserialize(serialized{0x40, 0xd0, 0x18, 0xfa})

	n.decafSqr(x)

	c.Assert(n, DeepEquals, m)
}

func (s *DecafSuite) Test_DecafMulW(c *C) {
	n := &bigNumber{}

	x, _ := deserialize(serialized{0xf5, 0x81, 0x74})
	y := int64(12363892)

	m, _ := deserialize(serialized{0x4, 0xab, 0xff, 0x19, 0xdc, 0x55})

	n.decafMulW(x, y)

	c.Assert(n, DeepEquals, m)
}

func (s *DecafSuite) Test_DecafSub(c *C) {
	x, _ := deserialize(serialized{
		0xf5, 0x81, 0x74, 0xd5, 0x7a, 0x33, 0x72,
		0x36, 0x3c, 0x0d, 0x9f, 0xcf, 0xaa, 0x3d,
		0xc1, 0x8b, 0x1e, 0xff, 0x7e, 0x89, 0xbf,
		0x76, 0x78, 0x63, 0x65, 0x80, 0xd1, 0x7d,
		0xd8, 0x4a, 0x87, 0x3b, 0x14, 0xb9, 0xc0,
		0xe1, 0x68, 0x0b, 0xbd, 0xc8, 0x76, 0x47,
		0xf3, 0xc3, 0x82, 0x90, 0x2d, 0x2f, 0x58,
		0xd2, 0x75, 0x4b, 0x39, 0xbc, 0xa8, 0x74,
	})

	z := &bigNumber{
		0xfffffff, 0xfffffff, 0xfffffff,
		0xfffffff, 0xfffffff, 0xfffffff,
		0xfffffff, 0xfffffff, 0xffffffe,
		0xfffffff, 0xfffffff, 0xfffffff,
		0xfffffff, 0xfffffff, 0xfffffff,
		0xfffffff}

	n := &bigNumber{}
	n.decafSub(x, x)

	c.Assert(n, DeepEquals, z)
}

func (s *Ed448Suite) Test_DecafCanon(c *C) {
	p, _ := deserialize(serialized{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})

	//p = p mod p = 0
	p.decafCanon()

	c.Assert(p, DeepEquals, &bigNumber{})
}

func (s *Ed448Suite) Test_DecafEq(c *C) {

	x, _ := deserialize(serialized{
		0xf5, 0x81, 0x74, 0xd5, 0x7a, 0x33, 0x72,
		0x36, 0x3c, 0x0d, 0x9f, 0xcf, 0xaa, 0x3d,
		0xc1, 0x8b, 0x1e, 0xff, 0x7e, 0x89, 0xbf,
		0x76, 0x78, 0x63, 0x65, 0x80, 0xd1, 0x7d,
		0xd8, 0x4a, 0x87, 0x3b, 0x14, 0xb9, 0xc0,
		0xe1, 0x68, 0x0b, 0xbd, 0xc8, 0x76, 0x47,
		0xf3, 0xc3, 0x82, 0x90, 0x2d, 0x2f, 0x58,
		0xd2, 0x75, 0x4b, 0x39, 0xbc, 0xa8, 0x74,
	})

	y, _ := deserialize(serialized{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})

	c.Assert(decafEq(x, x), Equals, true)
	c.Assert(decafEq(x, y), Equals, false)
}

func (s *Ed448Suite) Test_DecafAdd(c *C) {

	n := &bigNumber{}

	x, _ := deserialize(serialized{0x01})
	y, _ := deserialize(serialized{0x01})

	m, _ := deserialize(serialized{0x02})

	n.decafAdd(x, y)

	c.Assert(n, DeepEquals, m)
}

func (s *Ed448Suite) Test_DecafIsqrt(c *C) {

	y := mustDeserialize(serialized{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	})

	y.decafIsqrt(y)

	bs, _ := hex.DecodeString("04027d13a34bbe052fdf4247b02a4a3406268203a09076e56dee9dc2b699c4abc66f2832a677dfd0bf7e70ee72f01db170839717d1c64f02")
	exp := new(bigNumber).setBytes(bs)

	c.Assert(decafEq(y, exp), Equals, true)
}

func (s *Ed448Suite) Test_DecafConditionalNegateNumber(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	n := new(bigNumber).setBytes(bs)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	negated := new(bigNumber).setBytes(bs)

	x := &bigNumber{}
	x.copyFrom(n)
	x.decafCondNegate(0xffffffff)

	c.Assert(x, DeepEquals, negated)

	y := &bigNumber{}
	y.copyFrom(n)

	y.decafCondNegate(0)

	c.Assert(y, DeepEquals, n)
}
