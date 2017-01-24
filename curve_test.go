package ed448

import (
	"bytes"
	"encoding/hex"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestRadixBasePointIsOnCurve(c *C) {
	c.Assert(basePoint.OnCurve(), Equals, true)
}

func (s *Ed448Suite) TestRadixMultiplyByBase(c *C) {
	scalar := [scalarWords]word_t{}
	scalar[scalarWords-1] = 1000

	p := curve.multiplyByBase(scalar)

	c.Assert(p.OnCurve(), Equals, true)
}

func (s *Ed448Suite) TestRadixGenerateKey(c *C) {
	buffer := make([]byte, symKeyBytes)
	buffer[0] = 0x10
	r := bytes.NewReader(buffer[:])

	privKey, err := curve.generateKey(r)

	expectedSymKey := make([]byte, symKeyBytes)
	expectedSymKey[0] = 0x10

	expectedPriv := []byte{
		0xb3, 0xa4, 0x53, 0x31, 0xb1, 0x2b, 0x41, 0x1a,
		0xda, 0x51, 0xcf, 0xba, 0x0d, 0xea, 0x65, 0xb3,
		0x1b, 0x97, 0x9b, 0x41, 0xfe, 0x18, 0x93, 0x0c,
		0x6e, 0x4c, 0x02, 0x8a, 0x26, 0x24, 0xdf, 0xf0,
		0x24, 0x24, 0x06, 0x01, 0x4a, 0xb6, 0x3c, 0xab,
		0x33, 0x1e, 0xb5, 0xcf, 0x79, 0xc2, 0xc2, 0x6b,
		0xbb, 0x5e, 0xf8, 0xd8, 0x3e, 0x74, 0x26, 0x2c,
	}

	expectedPublic := []byte{
		0x73, 0xd7, 0xc8, 0x4a, 0x7e, 0x1d, 0x51, 0xea,
		0x17, 0x6c, 0xfb, 0x81, 0x31, 0x73, 0xcb, 0xba,
		0xf7, 0x20, 0x7c, 0xc0, 0x62, 0x1d, 0x80, 0xdc,
		0x60, 0x39, 0x64, 0x05, 0x4f, 0x12, 0xad, 0xaf,
		0x68, 0xf9, 0x1e, 0xcb, 0xaa, 0x38, 0xb6, 0xba,
		0x23, 0x96, 0x89, 0x67, 0xb6, 0x11, 0x6f, 0xb0,
		0xf8, 0x6e, 0x41, 0x92, 0x03, 0x67, 0x80, 0xdc,
	}

	c.Assert(err, IsNil)
	c.Assert(privKey.symKey(), DeepEquals, expectedSymKey)
	c.Assert(privKey.secretKey(), DeepEquals, expectedPriv)
	c.Assert(privKey.publicKey(), DeepEquals, expectedPublic)
}

func (s *Ed448Suite) TestDeriveNonce(c *C) {
	msg := []byte("hey there")
	symKey := [symKeyBytes]byte{
		0x27, 0x54, 0xcd, 0xa7, 0x12, 0x98, 0x88, 0x3d,
		0x4e, 0xf5, 0x11, 0x23, 0x92, 0x74, 0xb8, 0xa7,
		0xef, 0x7e, 0x51, 0x7e, 0x31, 0x28, 0xd4, 0xf7,
		0xfc, 0xfd, 0x9c, 0x62, 0xff, 0x65, 0x09, 0x65,
	}

	expectedNonce := [fieldWords]word_t{
		0xa7f2a3b8, 0xd4506099,
		0xfadfcc8b, 0xa1cdb278,
		0xe1228e40, 0x1a7b2b8c,
		0x4c0d9395, 0x01eeb2ac,
		0xd11846d5, 0x662e3ac6,
		0x8010aa06, 0x28ff671c,
		0xee92ab72, 0x066d13cb,
	}

	nonce := deriveNonce(msg, symKey[:])

	c.Assert(nonce, DeepEquals, expectedNonce)
}

func (s *Ed448Suite) TestDeriveChallenge(c *C) {
	msg := []byte("hey there")
	pubKey := [pubKeyBytes]byte{
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,
	}
	tmpSignature := [fieldBytes]uint8{
		0x66, 0x86, 0x04, 0xa8, 0x71, 0x4c, 0x39, 0xb9,
		0x42, 0x01, 0x7b, 0x45, 0xb6, 0xc7, 0xaf, 0xdb,
		0x7c, 0xad, 0x1f, 0x80, 0xa0, 0x23, 0x4d, 0xb5,
		0xab, 0x7c, 0x55, 0xf4, 0x38, 0x7d, 0xab, 0x60,
		0x25, 0x5a, 0x3d, 0xc9, 0xa1, 0x89, 0x85, 0xd1,
		0xc7, 0x4b, 0x19, 0x39, 0xbb, 0x08, 0x49, 0x09,
		0x0e, 0x0a, 0x31, 0x5a, 0x05, 0x5d, 0xe6, 0x47,
	}

	expectedChallenge := [fieldWords]word_t{
		0xee6472a4, 0x0007ea9a,
		0x58390769, 0xed3a8792,
		0x5fcb3d28, 0x38f60426,
		0x193bc68e, 0x7243ab24,
		0x2a4b0bd1, 0xa6c8365d,
		0x1a0f94f6, 0xab81dc7d,
		0xd46da3f1, 0x2db93e78,
	}

	challenge := deriveChallenge(pubKey[:], tmpSignature, msg)

	c.Assert(challenge, DeepEquals, expectedChallenge)
}

func (s *Ed448Suite) TestSign(c *C) {
	msg := []byte("hey there")
	k := privateKey([privKeyBytes]byte{
		//secret
		0x1f, 0x44, 0xfd, 0x2e, 0xde, 0x47, 0xca, 0xa8,
		0x7c, 0x4c, 0x45, 0x88, 0x1a, 0x7e, 0x01, 0x5a,
		0xa9, 0x01, 0x37, 0xfb, 0x0d, 0xbe, 0xb9, 0xe0,
		0xeb, 0x47, 0x29, 0xf7, 0x74, 0x0b, 0x5c, 0x23,
		0x66, 0xaa, 0xfd, 0x39, 0x03, 0x38, 0x78, 0x80,
		0x8f, 0xb2, 0x06, 0x13, 0x4e, 0xfb, 0xcf, 0x02,
		0x11, 0x43, 0x11, 0x3a, 0xd1, 0xf8, 0xb8, 0x22,

		//public
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,

		//symmetric
		0x27, 0x54, 0xcd, 0xa7, 0x12, 0x98, 0x88, 0x3d,
		0x4e, 0xf5, 0x11, 0x23, 0x92, 0x74, 0xb8, 0xa7,
		0xef, 0x7e, 0x51, 0x7e, 0x31, 0x28, 0xd4, 0xf7,
		0xfc, 0xfd, 0x9c, 0x62, 0xff, 0x65, 0x09, 0x65,
	})
	expectedSignature := [signatureBytes]byte{
		0xc7, 0x5d, 0xd7, 0x65, 0x5d, 0xfd, 0x70, 0xff,
		0x05, 0xa1, 0xe2, 0x02, 0x2d, 0x5c, 0x50, 0xaf,
		0x46, 0x7a, 0xa4, 0xea, 0x31, 0x51, 0x76, 0x8d,
		0xd0, 0x1d, 0x8c, 0x7a, 0x46, 0xb8, 0x08, 0x0b,
		0xdc, 0x06, 0x2a, 0xe6, 0xbf, 0x0e, 0x4a, 0x44,
		0x8a, 0x76, 0xdd, 0x52, 0xd3, 0x4f, 0x87, 0xe6,
		0x95, 0x36, 0x4b, 0xd2, 0x89, 0xdd, 0xcc, 0x82,
		0x49, 0xf6, 0xc4, 0x0e, 0x4c, 0x16, 0x18, 0x18,
		0xfd, 0x27, 0xe8, 0xeb, 0x3e, 0x81, 0x50, 0xe8,
		0x47, 0x9c, 0xa6, 0x99, 0x1d, 0x43, 0x0b, 0x53,
		0x22, 0xa6, 0xf1, 0x75, 0x8d, 0x7a, 0xec, 0x59,
		0xa9, 0xa4, 0xad, 0x92, 0x21, 0x84, 0x72, 0x00,
		0xfc, 0x5a, 0xa2, 0x4e, 0x05, 0x0c, 0x9b, 0x0a,
		0x8c, 0x3f, 0xc8, 0x46, 0x73, 0x42, 0x19, 0x35,
	}

	signature, err := curve.sign(msg, &k)

	c.Assert(err, IsNil)
	c.Assert(signature, DeepEquals, expectedSignature)
}

func (s *Ed448Suite) TestVerify(c *C) {
	msg := []byte("hey there")
	k := publicKey([pubKeyBytes]byte{
		//public
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,
	})
	signature := [signatureBytes]byte{
		0xc7, 0x5d, 0xd7, 0x65, 0x5d, 0xfd, 0x70, 0xff,
		0x05, 0xa1, 0xe2, 0x02, 0x2d, 0x5c, 0x50, 0xaf,
		0x46, 0x7a, 0xa4, 0xea, 0x31, 0x51, 0x76, 0x8d,
		0xd0, 0x1d, 0x8c, 0x7a, 0x46, 0xb8, 0x08, 0x0b,
		0xdc, 0x06, 0x2a, 0xe6, 0xbf, 0x0e, 0x4a, 0x44,
		0x8a, 0x76, 0xdd, 0x52, 0xd3, 0x4f, 0x87, 0xe6,
		0x95, 0x36, 0x4b, 0xd2, 0x89, 0xdd, 0xcc, 0x82,
		0x49, 0xf6, 0xc4, 0x0e, 0x4c, 0x16, 0x18, 0x18,
		0xfd, 0x27, 0xe8, 0xeb, 0x3e, 0x81, 0x50, 0xe8,
		0x47, 0x9c, 0xa6, 0x99, 0x1d, 0x43, 0x0b, 0x53,
		0x22, 0xa6, 0xf1, 0x75, 0x8d, 0x7a, 0xec, 0x59,
		0xa9, 0xa4, 0xad, 0x92, 0x21, 0x84, 0x72, 0x00,
		0xfc, 0x5a, 0xa2, 0x4e, 0x05, 0x0c, 0x9b, 0x0a,
		0x8c, 0x3f, 0xc8, 0x46, 0x73, 0x42, 0x19, 0x35,
	}

	valid := curve.verify(signature, msg, &k)

	c.Assert(valid, Equals, true)
}

func (s *Ed448Suite) TestMultiplyMontgomery(c *C) {
	pk := mustDeserialize(serialized{
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,
	})

	sk := [fieldWords]word_t{
		0x2efd441f, 0xa8ca47de,
		0x88454c7c, 0x5a017e1a,
		0xfb3701a9, 0xe0b9be0d,
		0xf72947eb, 0x235c0b74,
		0x39fdaa66, 0x80783803,
		0x1306b28f, 0x02cffb4e,
		0x3a114311, 0x22b8f8d1,
	}

	bs, _ := hex.DecodeString("322d71661943b5e080abed64d9ed331874a975329aaf9b42815e793ac08691e478fe559b29593a5413d5a4475e3ae0735a6d9bc1dc192b7d")
	expectedPublic := new(bigNumber)
	expectedPublic.setBytes(bs)

	pk, ok := curve.multiplyMontgomery(pk, sk, scalarBits, 1)

	c.Assert(ok, Equals, word_t(0))
	c.Assert(pk, DeepEquals, expectedPublic)
}

func (s *Ed448Suite) Test_DecafDerivePrivate(c *C) {
	sym := [symKeyBytes]byte{
		0x1f, 0x16, 0x6c, 0x08, 0xc7, 0xc1, 0x41, 0xb0,
		0x49, 0xa2, 0x80, 0x3a, 0xcf, 0x4a, 0x82, 0x84,
		0x13, 0x4f, 0x7c, 0x72, 0x89, 0xa1, 0x1d, 0xc5,
		0xa6, 0x0e, 0x0c, 0xc2, 0x7b, 0x9c, 0xbb, 0x87,
	}

	pub := []byte{
		0x3d, 0x28, 0x7d, 0x7f, 0xcc, 0xad, 0x34, 0x32,
		0x66, 0x1f, 0x6d, 0x9a, 0x9b, 0xc3, 0x15, 0xe4,
		0x07, 0x5e, 0x9d, 0x02, 0xe2, 0xd2, 0x2f, 0xfe,
		0x20, 0x5d, 0x0b, 0x28, 0xda, 0x0c, 0x4b, 0x96,
		0x41, 0x83, 0x33, 0x11, 0xa4, 0x93, 0x0c, 0x27,
		0xa1, 0xaf, 0xce, 0x49, 0x41, 0xfc, 0x57, 0x17,
		0x38, 0x9d, 0x96, 0x6f, 0x91, 0x51, 0x19, 0x25,
	}

	priv := []byte{
		0xcb, 0x6a, 0xc3, 0x11, 0x80, 0x8c, 0xdf, 0xd6,
		0xc4, 0x9b, 0x0b, 0x9b, 0x2b, 0x96, 0x1e, 0x92,
		0x33, 0x26, 0xe5, 0x6f, 0xbd, 0x5a, 0x68, 0xec,
		0x95, 0x8b, 0x89, 0x87, 0xda, 0x12, 0x8d, 0x44,
		0xf0, 0x2a, 0x24, 0xe4, 0x3a, 0xf3, 0x77, 0xc8,
		0xa7, 0xba, 0xee, 0xa0, 0x6e, 0xf6, 0x7b, 0xf0,
		0x07, 0x22, 0xa6, 0x4a, 0xeb, 0x8e, 0xce, 0x02,
	}

	pk, _ := curve.decafDerivePrivateKey(sym)

	sym2 := []byte{
		0x1f, 0x16, 0x6c, 0x08, 0xc7, 0xc1, 0x41, 0xb0,
		0x49, 0xa2, 0x80, 0x3a, 0xcf, 0x4a, 0x82, 0x84,
		0x13, 0x4f, 0x7c, 0x72, 0x89, 0xa1, 0x1d, 0xc5,
		0xa6, 0x0e, 0x0c, 0xc2, 0x7b, 0x9c, 0xbb, 0x87,
	}

	c.Assert(pk.symKey(), DeepEquals, sym2)
	c.Assert(pk.secretKey(), DeepEquals, priv)
	c.Assert(pk.publicKey(), DeepEquals, pub)
}
