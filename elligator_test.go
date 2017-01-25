package ed448

import (
	. "gopkg.in/check.v1"
)

func (s *DecafSuite) Test_2Torsion(c *C) {

	p := &pointT{
		&bigNumber{0x00},
		&bigNumber{0x01},
		&bigNumber{0x01},
		&bigNumber{0x00},
	}

	q := &pointT{
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
	}

	e := &pointT{
		// gives you p
		&bigNumber{0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xffffffe, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff},
		&bigNumber{0xffffffe, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xffffffe, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff, 0xfffffff},
		&bigNumber{0x01},
		&bigNumber{0x00},
	}

	point2Torque(q, p)

	c.Assert(q.x, DeepEquals, e.x)
	c.Assert(q.y, DeepEquals, e.y)
	c.Assert(q.z, DeepEquals, e.z)
	c.Assert(q.t, DeepEquals, e.t)

}
