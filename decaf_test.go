package ed448

import (
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

	// this is not the same result as the one in Karatzuba
	z := &bigNumber{
		0xaf1b9c5, 0xe417cd7, 0x839a472, 0x43bcebc,
		0x7dcf7e0, 0x726193c, 0x304a5bb, 0xf04e22e,
		0x526f560, 0x44604e2, 0x54f3a45, 0x979c291,
		0x57bc577, 0x431d314, 0xf5bf68e, 0x64baf6ec}

	result := new(bigNumber)

	ser := decafMul(result, x, y)

	c.Assert(ser, DeepEquals, z)
}
