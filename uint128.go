// Very simple implementation of uint128
package numeric

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Uint128 struct {
	High uint64
	Low  uint64
}

func (u *Uint128) Bytes() []byte {
	bytes := make([]byte, 16)
	binary.LittleEndian.PutUint64(bytes[0:8], u.Low)
	binary.LittleEndian.PutUint64(bytes[8:16], u.High)
	return bytes
}

func Uint128FromLittleEndianBytes(bytes []byte) Uint128 {
	return Uint128{
		High: binary.LittleEndian.Uint64(bytes[8:16]),
		Low:  binary.LittleEndian.Uint64(bytes[0:8]),
	}
}

func Uint128FromBigEndianBytes(bytes []byte) Uint128 {
	return Uint128{
		High: binary.BigEndian.Uint64(bytes[0:8]),
		Low:  binary.BigEndian.Uint64(bytes[8:16]),
	}
}

func (u Uint128) FormatString() string {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes[0:8], u.High)
	binary.BigEndian.PutUint64(bytes[8:16], u.Low)
	return FormatBigEndianBinaryString(bytes, "_")
}

func (u Uint128) String() string {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes[0:8], u.High)
	binary.BigEndian.PutUint64(bytes[8:16], u.Low)
	return FormatBigEndianBinaryString(bytes, "")
}

func (u Uint128) And(other Uint128) Uint128 {
	return Uint128{
		High: u.High & other.High,
		Low:  u.Low & other.Low,
	}
}

func (u Uint128) Or(other Uint128) Uint128 {
	return Uint128{
		High: u.High | other.High,
		Low:  u.Low | other.Low,
	}
}

func (u Uint128) Xor(other Uint128) Uint128 {
	return Uint128{
		High: u.High ^ other.High,
		Low:  u.Low ^ other.Low,
	}
}

// Mask returns a mask of the given number of bits
// Acquire bits less equals than 128
func Mask(bits uint8) Uint128 {
	if bits == 0 {
		return Uint128{
			High: 0,
			Low:  0,
		}
	}
	if bits>>7 != 0 {
		return Uint128{
			High: ^uint64(0),
			Low:  ^uint64(0),
		}
	}
	// bits < 0x7F(128)
	if bits&0x40 != 0 { // greater equal than 64
		return Uint128{
			High: MaskUint64(bits & 0x3F),
			Low:  ^uint64(0),
		}
	}
	return Uint128{
		High: 0,
		Low:  MaskUint64(bits), // equals to min(bits, 64)
	}
}

func MaskUint64(bits uint8) uint64 {
	return ^uint64(0) >> (64 - bits)
}

func (u Uint128) ShiftRight(shiftAmount uint8) Uint128 {

	if shiftAmount == 0 {
		return u
	}
	if shiftAmount>>7 != 0 {
		return Uint128{
			High: 0,
			Low:  0,
		}
	}
	// shiftAmount < 0x7F(128)
	if shiftAmount&0x40 != 0 { // equals to shiftAmount >= 64
		return Uint128{
			High: 0,
			Low:  u.High >> (shiftAmount & 0x3F), // equals to shiftAmount - 64
		}
	}
	return Uint128{
		High: u.High >> shiftAmount,
		Low:  u.Low>>shiftAmount | u.High<<(64-shiftAmount),
	}
}

func (u Uint128) ShiftLeft(shiftAmount uint8) Uint128 {
	if shiftAmount == 0 {
		return u
	}
	if shiftAmount >= 128 {
		return Uint128{
			High: 0,
			Low:  0,
		}
	}
	// shiftAmount < 0x7F(128)
	if shiftAmount&0x40 != 0 { // equals to shiftAmount >= 64
		return Uint128{
			High: u.Low << (shiftAmount & 0x3F), // equals to shiftAmount - 64
			Low:  0,
		}
	}
	return Uint128{
		High: u.High<<shiftAmount | u.Low>>(64-shiftAmount),
		Low:  u.Low << shiftAmount,
	}
}

func SwapBytes(bytes []byte) []byte {
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i]
	}
	return bytes
}

func SwapBytes16(bytes []byte) []byte {
	bytes[0], bytes[15] = bytes[15], bytes[0]
	bytes[1], bytes[14] = bytes[14], bytes[1]
	bytes[2], bytes[13] = bytes[13], bytes[2]
	bytes[3], bytes[12] = bytes[12], bytes[3]
	bytes[4], bytes[11] = bytes[11], bytes[4]
	bytes[5], bytes[10] = bytes[10], bytes[5]
	bytes[6], bytes[9] = bytes[9], bytes[6]
	bytes[7], bytes[8] = bytes[8], bytes[7]
	return bytes
}

func FormatBigEndianBinaryString(bigEndianBytes []byte, sep string) string {
	var builder strings.Builder
	for i := 0; i < len(bigEndianBytes); i++ {
		if i != 0 {
			builder.WriteString(sep)
		}
		b := bigEndianBytes[i]
		c := fmt.Sprintf("%08b", b)
		builder.WriteString(c)
	}
	return builder.String()
}
