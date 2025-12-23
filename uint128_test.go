package numeric_test

import (
	"testing"

	numeric "github.com/everoute/numeric-go"
	. "github.com/onsi/gomega"
)

func TestMark0(t *testing.T) {
	RegisterTestingT(t)
	mask := numeric.Mask(0)
	str := mask.FormatString()
	Expect(str).Should(HaveLen(128 + 15))
	Expect(str[0 : 64+7]).Should(Equal("00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000"))
	Expect(str[64+8:]).Should(Equal("00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000"))
}

func TestMark64(t *testing.T) {
	RegisterTestingT(t)
	mask := numeric.Mask(64)
	str := mask.FormatString()
	Expect(str).Should(HaveLen(128 + 15))
	Expect(str[0 : 64+7]).Should(Equal("00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000"))
	Expect(str[64+8:]).Should(Equal("11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111"))
}

func TestMark128(t *testing.T) {
	RegisterTestingT(t)
	mask := numeric.Mask(128)
	str := mask.FormatString()
	Expect(str).Should(HaveLen(128 + 15))
	Expect(str[0 : 64+7]).Should(Equal("11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111"))
	Expect(str[64+8:]).Should(Equal("11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111"))
}

func TestMaskLessThan64(t *testing.T) {
	RegisterTestingT(t)
	// mask := numeric.OriginInportMask
	mask := numeric.Mask(16)
	str := mask.FormatString()
	Expect(str).Should(HaveLen(128 + 15))
	Expect(str[0 : 64+7]).Should(Equal("00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000"))
	Expect(str[64+8:]).Should(Equal("00000000_00000000_00000000_00000000_00000000_00000000_11111111_11111111"))
}

func TestMaskGreaterThan64(t *testing.T) {
	RegisterTestingT(t)
	mask := numeric.Mask(96)
	str := mask.FormatString()
	Expect(str).Should(HaveLen(128 + 15))
	Expect(str[0 : 64+7]).Should(Equal("00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111"))
	Expect(str[64+8:]).Should(Equal("11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111"))
}

func TestShiftLeft(t *testing.T) {
	RegisterTestingT(t)
	type testCase struct {
		old      numeric.Uint128
		shift    uint8
		expected numeric.Uint128
	}
	cases := []testCase{
		// bits less than 64

		// shift zero
		{numeric.Uint128{High: 0, Low: numeric.MaskUint64(16)}, 0, numeric.Uint128{High: 0, Low: numeric.MaskUint64(16)}},
		// shift 64
		{numeric.Uint128{High: 0, Low: numeric.MaskUint64(16)}, 64, numeric.Uint128{High: numeric.MaskUint64(16), Low: 0}},
		// shift 128
		{numeric.Uint128{High: 0, Low: numeric.MaskUint64(16)}, 128, numeric.Uint128{High: 0, Low: 0}},
		// shift less than 64
		{numeric.Uint128{High: 0, Low: numeric.MaskUint64(16)}, 32, numeric.Uint128{High: 0, Low: numeric.MaskUint64(16) << 32}},
		// shift greater than 64
		{numeric.Uint128{High: 0, Low: numeric.MaskUint64(16)}, 96, numeric.Uint128{High: numeric.MaskUint64(16) << (96 - 64), Low: 0}},

		// bits greater than 64

		// shift zero
		{numeric.Uint128{High: numeric.MaskUint64(16), Low: ^uint64(0)}, 0, numeric.Uint128{High: numeric.MaskUint64(16), Low: ^uint64(0)}},
		// shift 64
		{numeric.Uint128{High: numeric.MaskUint64(16), Low: ^uint64(0)}, 64, numeric.Uint128{High: ^uint64(0), Low: 0}},
		// shift 128
		{numeric.Uint128{High: numeric.MaskUint64(16), Low: ^uint64(0)}, 128, numeric.Uint128{High: 0, Low: 0}},
		// shift less than 64
		{numeric.Uint128{High: numeric.MaskUint64(16), Low: ^uint64(0)}, 32, numeric.Uint128{High: numeric.MaskUint64(48), Low: numeric.MaskUint64(32) << 32}},
		// shift greater than 64
		{numeric.Uint128{High: numeric.MaskUint64(16), Low: ^uint64(0)}, 96, numeric.Uint128{High: numeric.MaskUint64(32) << 32, Low: 0}},
	}

	for _, c := range cases {
		new := c.old.ShiftLeft(c.shift)
		c1 := testCase{c.old, c.shift, new}
		Expect(c1).Should(Equal(c))
	}
}

func TestShiftRight(t *testing.T) {
	RegisterTestingT(t)
	type testCase struct {
		old      numeric.Uint128
		shift    uint8
		expected numeric.Uint128
	}
	cases := []testCase{
		// bits less than 64

		// shift zero
		{numeric.Uint128{High: ^numeric.MaskUint64(48), Low: 0}, 0, numeric.Uint128{High: ^numeric.MaskUint64(48), Low: 0}},
		// shift 64
		{numeric.Uint128{High: ^numeric.MaskUint64(48), Low: 0}, 64, numeric.Uint128{High: 0, Low: ^numeric.MaskUint64(48)}},
		// shift 128
		{numeric.Uint128{High: ^numeric.MaskUint64(48), Low: 0}, 128, numeric.Uint128{High: 0, Low: 0}},
		// shift less than 64
		{numeric.Uint128{High: ^numeric.MaskUint64(48), Low: 0}, 32, numeric.Uint128{High: ^numeric.MaskUint64(48) >> 32, Low: 0}},
		// shift greater than 64
		{numeric.Uint128{High: ^numeric.MaskUint64(48), Low: 0}, 96, numeric.Uint128{High: 0, Low: ^numeric.MaskUint64(48) >> (96 - 64)}},

		// bits greater than 64

		// shift zero
		{numeric.Uint128{High: ^uint64(0), Low: ^numeric.MaskUint64(48)}, 0, numeric.Uint128{High: ^uint64(0), Low: ^numeric.MaskUint64(48)}},
		// shift 64
		{numeric.Uint128{High: ^uint64(0), Low: ^numeric.MaskUint64(48)}, 64, numeric.Uint128{High: 0, Low: ^uint64(0)}},
		// shift 128
		{numeric.Uint128{High: ^uint64(0), Low: ^numeric.MaskUint64(48)}, 128, numeric.Uint128{High: 0, Low: 0}},
		// shift less than 64
		{numeric.Uint128{High: ^uint64(0), Low: ^numeric.MaskUint64(48)}, 32, numeric.Uint128{High: numeric.MaskUint64(32), Low: ^numeric.MaskUint64(16)}},
		// shift greater than 64
		{numeric.Uint128{High: ^uint64(0), Low: ^numeric.MaskUint64(48)}, 96, numeric.Uint128{High: 0, Low: numeric.MaskUint64(32)}},
	}

	for _, c := range cases {
		new := c.old.ShiftRight(c.shift)
		c1 := testCase{c.old, c.shift, new}
		Expect(c1).Should(Equal(c))
	}
}
