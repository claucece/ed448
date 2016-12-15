package ed448

const (
	// D is the non-square element of F_p
	d = -39081

	montgomeryFactor = "3bd440fae918bc5ull"
)

// XXX: Will bigNumber change with Decaf?
//type decaf_word_t uint32

//type gf struct {
//	limb [16]decaf_word_t
//}

// Copy copies n = y
func (n *bigNumber) decafCopy(x *bigNumber) *bigNumber {
	copy(n[:], x[:])
	return n
}
