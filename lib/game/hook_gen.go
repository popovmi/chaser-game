package game

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Hook) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "end":
			err = z.End.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "End")
				return
			}
		case "vel":
			err = z.Vel.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Vel")
				return
			}
		case "current_distance":
			z.CurrentDistance, err = dc.ReadFloat64()
			if err != nil {
				err = msgp.WrapError(err, "CurrentDistance")
				return
			}
		case "stuck":
			z.Stuck, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Stuck")
				return
			}
		case "is_returning":
			z.IsReturning, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "IsReturning")
				return
			}
		case "caught_player_id":
			z.CaughtPlayerID, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "CaughtPlayerID")
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
func (z *Hook) EncodeMsg(en *msgp.Writer) (err error) {
	// check for omitted fields
	zb0001Len := uint32(6)
	var zb0001Mask uint8 /* 6 bits */
	_ = zb0001Mask
	if z.CurrentDistance == 0 {
		zb0001Len--
		zb0001Mask |= 0x4
	}
	if z.Stuck == false {
		zb0001Len--
		zb0001Mask |= 0x8
	}
	if z.IsReturning == false {
		zb0001Len--
		zb0001Mask |= 0x10
	}
	if z.CaughtPlayerID == "" {
		zb0001Len--
		zb0001Mask |= 0x20
	}
	// variable map header, size zb0001Len
	err = en.Append(0x80 | uint8(zb0001Len))
	if err != nil {
		return
	}

	// skip if no fields are to be emitted
	if zb0001Len != 0 {
		// write "end"
		err = en.Append(0xa3, 0x65, 0x6e, 0x64)
		if err != nil {
			return
		}
		err = z.End.EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "End")
			return
		}
		// write "vel"
		err = en.Append(0xa3, 0x76, 0x65, 0x6c)
		if err != nil {
			return
		}
		err = z.Vel.EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "Vel")
			return
		}
		if (zb0001Mask & 0x4) == 0 { // if not omitted
			// write "current_distance"
			err = en.Append(0xb0, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65)
			if err != nil {
				return
			}
			err = en.WriteFloat64(z.CurrentDistance)
			if err != nil {
				err = msgp.WrapError(err, "CurrentDistance")
				return
			}
		}
		if (zb0001Mask & 0x8) == 0 { // if not omitted
			// write "stuck"
			err = en.Append(0xa5, 0x73, 0x74, 0x75, 0x63, 0x6b)
			if err != nil {
				return
			}
			err = en.WriteBool(z.Stuck)
			if err != nil {
				err = msgp.WrapError(err, "Stuck")
				return
			}
		}
		if (zb0001Mask & 0x10) == 0 { // if not omitted
			// write "is_returning"
			err = en.Append(0xac, 0x69, 0x73, 0x5f, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x69, 0x6e, 0x67)
			if err != nil {
				return
			}
			err = en.WriteBool(z.IsReturning)
			if err != nil {
				err = msgp.WrapError(err, "IsReturning")
				return
			}
		}
		if (zb0001Mask & 0x20) == 0 { // if not omitted
			// write "caught_player_id"
			err = en.Append(0xb0, 0x63, 0x61, 0x75, 0x67, 0x68, 0x74, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64)
			if err != nil {
				return
			}
			err = en.WriteString(z.CaughtPlayerID)
			if err != nil {
				err = msgp.WrapError(err, "CaughtPlayerID")
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Hook) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// check for omitted fields
	zb0001Len := uint32(6)
	var zb0001Mask uint8 /* 6 bits */
	_ = zb0001Mask
	if z.CurrentDistance == 0 {
		zb0001Len--
		zb0001Mask |= 0x4
	}
	if z.Stuck == false {
		zb0001Len--
		zb0001Mask |= 0x8
	}
	if z.IsReturning == false {
		zb0001Len--
		zb0001Mask |= 0x10
	}
	if z.CaughtPlayerID == "" {
		zb0001Len--
		zb0001Mask |= 0x20
	}
	// variable map header, size zb0001Len
	o = append(o, 0x80|uint8(zb0001Len))

	// skip if no fields are to be emitted
	if zb0001Len != 0 {
		// string "end"
		o = append(o, 0xa3, 0x65, 0x6e, 0x64)
		o, err = z.End.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "End")
			return
		}
		// string "vel"
		o = append(o, 0xa3, 0x76, 0x65, 0x6c)
		o, err = z.Vel.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Vel")
			return
		}
		if (zb0001Mask & 0x4) == 0 { // if not omitted
			// string "current_distance"
			o = append(o, 0xb0, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65)
			o = msgp.AppendFloat64(o, z.CurrentDistance)
		}
		if (zb0001Mask & 0x8) == 0 { // if not omitted
			// string "stuck"
			o = append(o, 0xa5, 0x73, 0x74, 0x75, 0x63, 0x6b)
			o = msgp.AppendBool(o, z.Stuck)
		}
		if (zb0001Mask & 0x10) == 0 { // if not omitted
			// string "is_returning"
			o = append(o, 0xac, 0x69, 0x73, 0x5f, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x69, 0x6e, 0x67)
			o = msgp.AppendBool(o, z.IsReturning)
		}
		if (zb0001Mask & 0x20) == 0 { // if not omitted
			// string "caught_player_id"
			o = append(o, 0xb0, 0x63, 0x61, 0x75, 0x67, 0x68, 0x74, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64)
			o = msgp.AppendString(o, z.CaughtPlayerID)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Hook) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "end":
			bts, err = z.End.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "End")
				return
			}
		case "vel":
			bts, err = z.Vel.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Vel")
				return
			}
		case "current_distance":
			z.CurrentDistance, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "CurrentDistance")
				return
			}
		case "stuck":
			z.Stuck, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Stuck")
				return
			}
		case "is_returning":
			z.IsReturning, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "IsReturning")
				return
			}
		case "caught_player_id":
			z.CaughtPlayerID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "CaughtPlayerID")
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
func (z *Hook) Msgsize() (s int) {
	s = 1 + 4 + z.End.Msgsize() + 4 + z.Vel.Msgsize() + 17 + msgp.Float64Size + 6 + msgp.BoolSize + 13 + msgp.BoolSize + 17 + msgp.StringPrefixSize + len(z.CaughtPlayerID)
	return
}
