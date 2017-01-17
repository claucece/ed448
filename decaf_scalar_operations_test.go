package ed448

import . "gopkg.in/check.v1"

var (
	niels = &twNiels{
		&bigNumber{0x068d5b74},
		&bigNumber{0x068d5b74},
		&bigNumber{0x068d5b74},
	}

	oldProjective = &pointT{
		&bigNumber{},
		&bigNumber{0x00000001},
		&bigNumber{0x00000001},
		&bigNumber{},
	}
)

func resetValues() {
	oldProjective = &pointT{
		&bigNumber{},
		&bigNumber{0x00000001},
		&bigNumber{0x00000001},
		&bigNumber{},
	}
}

func resetPoint(p *pointT) {
	p = &pointT{
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
	}
}

func (s *Ed448Suite) TearDownTest(c *C) {
	resetValues()
}

func (s *Ed448Suite) Test_ScalarAdditionAndSubtraction(c *C) {

	scalar1 := [scalarWords]word_t{
		0x529eec33, 0x721cf5b5,
		0xc8e9c2ab, 0x7a4cf635,
		0x44a725bf, 0xeec492d9,
		0x0cd77058, 0x00000002,
	}

	scalar2 := [scalarWords]word_t{
		0x00000001,
	}

	subExp := [scalarWords]word_t{
		0x529eec32, 0x721cf5b5,
		0xc8e9c2ab, 0x7a4cf635,
		0x44a725bf, 0xeec492d9,
		0x0cd77058, 0x00000002,
	}

	addExp := [scalarWords]word_t{
		0x529eec34, 0x721cf5b5,
		0xc8e9c2ab, 0x7a4cf635,
		0x44a725bf, 0xeec492d9,
		0x0cd77058, 0x00000002,
	}

	added := scalarAdd(scalar1, scalar2)
	subtracted := scalarSub(scalar1, scalar2)

	c.Assert(added, DeepEquals, addExp)
	c.Assert(subtracted, DeepEquals, subExp)
}

func (s *Ed448Suite) Test_ScalarHalve(c *C) {

	scalar1 := [scalarWords]word_t{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12,
	}

	scalar2 := [scalarWords]word_t{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4,
	}

	expected := [scalarWords]word_t{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
	}

	halved := scHalve(scalar1, scalar2)

	c.Assert(halved, DeepEquals, expected)
}

func (s *Ed448Suite) Test_PointDouble(c *C) {

	q := &pointT{
		&bigNumber{0x00000001},
		&bigNumber{0x00000002},
		&bigNumber{0x00000003},
		&bigNumber{0x00000004},
	}

	expX := &bigNumber{0x0000003b, 0x10000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0ffffffe, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}

	expY := &bigNumber{0x0000000e, 0x00000000,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}

	expZ := &bigNumber{0x0000002c, 0x10000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0ffffffe, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}

	expT := &bigNumber{0x00000013, 0x00000000,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}

	expected := &pointT{
		expX,
		expY,
		expZ,
		expT}

	p := &pointT{
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
	}

	p.pointDoubleInternal(q, false)

	c.Assert(p, DeepEquals, expected)

	resetPoint(p)
	p.pointDoubleInternal(q, true)

	exp1X := &bigNumber{0x0000003b, 0x10000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0ffffffe, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}

	exp1Y := &bigNumber{0x0000000e, 0x00000000,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
		0x00000000, 0x00000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}

	exp1Z := &bigNumber{0x0000002c, 0x10000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0ffffffe, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}

	exp1T := &bigNumber{0x00000002, 0x10000000,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0ffffffe, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
		0x0fffffff, 0x0fffffff,
	}
	expected1 := &pointT{
		exp1X,
		exp1Y,
		exp1Z,
		exp1T}

	c.Assert(p, DeepEquals, expected1)

}

func (s *Ed448Suite) Test_GenerateConstant(c *C) {

	c.Skip("In progress")

	adjustmentConstant := [scalarWords]word_t{0x529eec33,
		0x721cf5b5,
		0xc8e9c2ab,
		0x7a4cf635,
		0x44a725bf,
		0xeec492d9,
		0x0cd77058,
		0x00000002,
		0x00000000,
		0x00000000,
		0x00000000,
		0x00000000,
		0x00000000,
		0x00000000,
	}

	sc := scalarAdjustment()

	c.Assert(sc, DeepEquals, adjustmentConstant)
}

func (s *Ed448Suite) Test_AddNielsToProjective(c *C) {
	expected := &pointT{
		&bigNumber{0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0ffffffe, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
		&bigNumber{0x0d1ab6e7, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
		&bigNumber{0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		}, &bigNumber{},
	}

	oldProjective.addNielsToProjective(niels, true)

	c.Assert(oldProjective.x, DeepEquals, expected.x)
	c.Assert(oldProjective.y, DeepEquals, expected.y)
	c.Assert(oldProjective.z, DeepEquals, expected.z)
	c.Assert(oldProjective.t, DeepEquals, expected.t)
}

func (s *Ed448Suite) Test_AddNielsToProjective_BeforeDouble(c *C) {
	expectedProjective := &pointT{
		&bigNumber{0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0ffffffe, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
		&bigNumber{0x0d1ab6e7, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
		&bigNumber{0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x00000000, 0x00000000,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
		&bigNumber{0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0ffffffe, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
	}

	oldProjective.addNielsToProjective(niels, false)

	c.Assert(oldProjective.x, DeepEquals, expectedProjective.x)
	c.Assert(oldProjective.y, DeepEquals, expectedProjective.y)
	c.Assert(oldProjective.z, DeepEquals, expectedProjective.z)
	c.Assert(oldProjective.t, DeepEquals, expectedProjective.t)
}

func (s *Ed448Suite) Test_ConvertNielsToProjective(c *C) {
	expected := &pointT{
		&bigNumber{0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0ffffffe, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
		&bigNumber{0x0d1ab6e8},
		&bigNumber{0x00000001},
		&bigNumber{0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0ffffffe, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
			0x0fffffff, 0x0fffffff,
		},
	}

	convertNielsToPt(oldProjective, niels)

	c.Assert(oldProjective.x, DeepEquals, expected.x)
	c.Assert(oldProjective.y, DeepEquals, expected.y)
	c.Assert(oldProjective.z, DeepEquals, expected.z)
	c.Assert(oldProjective.t, DeepEquals, expected.t)
}

func (s *Ed448Suite) Test_CondNegNiels(c *C) {

}
