package vector

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Vector2D) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "x":
			z.X, err = dc.ReadFloat64()
			if err != nil {
				err = msgp.WrapError(err, "X")
				return
			}
		case "y":
			z.Y, err = dc.ReadFloat64()
			if err != nil {
				err = msgp.WrapError(err, "Y")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Vector2D) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "x"
	err = en.Append(0x82, 0xa1, 0x78)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.X)
	if err != nil {
		err = msgp.WrapError(err, "X")
		return
	}
	// write "y"
	err = en.Append(0xa1, 0x79)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.Y)
	if err != nil {
		err = msgp.WrapError(err, "Y")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Vector2D) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "x"
	o = append(o, 0x82, 0xa1, 0x78)
	o = msgp.AppendFloat64(o, z.X)
	// string "y"
	o = append(o, 0xa1, 0x79)
	o = msgp.AppendFloat64(o, z.Y)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Vector2D) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "x":
			z.X, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "X")
				return
			}
		case "y":
			z.Y, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Y")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Vector2D) Msgsize() (s int) {
	s = 1 + 2 + msgp.Float64Size + 2 + msgp.Float64Size
	return
}
