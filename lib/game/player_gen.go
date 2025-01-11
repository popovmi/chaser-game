package game

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Player) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "id":
			z.ID, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "color":
			err = z.Color.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Color")
				return
			}
		case "joined_at":
			z.JoinedAt, err = dc.ReadTime()
			if err != nil {
				err = msgp.WrapError(err, "JoinedAt")
				return
			}
		case "position":
			err = z.Position.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Position")
				return
			}
		case "velocity":
			err = z.Velocity.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Velocity")
				return
			}
		case "angle":
			z.Angle, err = dc.ReadFloat64()
			if err != nil {
				err = msgp.WrapError(err, "Angle")
				return
			}
		case "move_dir":
			z.MoveDir, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "MoveDir")
				return
			}
		case "turn_dir":
			{
				var zb0002 byte
				zb0002, err = dc.ReadByte()
				if err != nil {
					err = msgp.WrapError(err, "RotationDir")
					return
				}
				z.RotationDir = RotationDirection(zb0002)
			}
		case "hook":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "Hook")
					return
				}
				z.Hook = nil
			} else {
				if z.Hook == nil {
					z.Hook = new(Hook)
				}
				err = z.Hook.DecodeMsg(dc)
				if err != nil {
					err = msgp.WrapError(err, "Hook")
					return
				}
			}
		case "hooked_at":
			z.HookedAt, err = dc.ReadTime()
			if err != nil {
				err = msgp.WrapError(err, "HookedAt")
				return
			}
		case "is_hooked":
			z.IsHooked, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "IsHooked")
				return
			}
		case "caught_by_id":
			z.CaughtByID, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "CaughtByID")
				return
			}
		case "blinking":
			z.Blinking, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Blinking")
				return
			}
		case "blinked_at":
			z.BlinkedAt, err = dc.ReadTime()
			if err != nil {
				err = msgp.WrapError(err, "BlinkedAt")
				return
			}
		case "blinked":
			z.Blinked, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Blinked")
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
func (z *Player) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 16
	// write "id"
	err = en.Append(0xde, 0x0, 0x10, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "name"
	err = en.Append(0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	// write "color"
	err = en.Append(0xa5, 0x63, 0x6f, 0x6c, 0x6f, 0x72)
	if err != nil {
		return
	}
	err = z.Color.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Color")
		return
	}
	// write "joined_at"
	err = en.Append(0xa9, 0x6a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x5f, 0x61, 0x74)
	if err != nil {
		return
	}
	err = en.WriteTime(z.JoinedAt)
	if err != nil {
		err = msgp.WrapError(err, "JoinedAt")
		return
	}
	// write "position"
	err = en.Append(0xa8, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = z.Position.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Position")
		return
	}
	// write "velocity"
	err = en.Append(0xa8, 0x76, 0x65, 0x6c, 0x6f, 0x63, 0x69, 0x74, 0x79)
	if err != nil {
		return
	}
	err = z.Velocity.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Velocity")
		return
	}
	// write "angle"
	err = en.Append(0xa5, 0x61, 0x6e, 0x67, 0x6c, 0x65)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.Angle)
	if err != nil {
		err = msgp.WrapError(err, "Angle")
		return
	}
	// write "move_dir"
	err = en.Append(0xa8, 0x6d, 0x6f, 0x76, 0x65, 0x5f, 0x64, 0x69, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.MoveDir)
	if err != nil {
		err = msgp.WrapError(err, "MoveDir")
		return
	}
	// write "turn_dir"
	err = en.Append(0xa8, 0x74, 0x75, 0x72, 0x6e, 0x5f, 0x64, 0x69, 0x72)
	if err != nil {
		return
	}
	err = en.WriteByte(byte(z.RotationDir))
	if err != nil {
		err = msgp.WrapError(err, "RotationDir")
		return
	}
	// write "hook"
	err = en.Append(0xa4, 0x68, 0x6f, 0x6f, 0x6b)
	if err != nil {
		return
	}
	if z.Hook == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Hook.EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "Hook")
			return
		}
	}
	// write "hooked_at"
	err = en.Append(0xa9, 0x68, 0x6f, 0x6f, 0x6b, 0x65, 0x64, 0x5f, 0x61, 0x74)
	if err != nil {
		return
	}
	err = en.WriteTime(z.HookedAt)
	if err != nil {
		err = msgp.WrapError(err, "HookedAt")
		return
	}
	// write "is_hooked"
	err = en.Append(0xa9, 0x69, 0x73, 0x5f, 0x68, 0x6f, 0x6f, 0x6b, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(z.IsHooked)
	if err != nil {
		err = msgp.WrapError(err, "IsHooked")
		return
	}
	// write "caught_by_id"
	err = en.Append(0xac, 0x63, 0x61, 0x75, 0x67, 0x68, 0x74, 0x5f, 0x62, 0x79, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.CaughtByID)
	if err != nil {
		err = msgp.WrapError(err, "CaughtByID")
		return
	}
	// write "blinking"
	err = en.Append(0xa8, 0x62, 0x6c, 0x69, 0x6e, 0x6b, 0x69, 0x6e, 0x67)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Blinking)
	if err != nil {
		err = msgp.WrapError(err, "Blinking")
		return
	}
	// write "blinked_at"
	err = en.Append(0xaa, 0x62, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x5f, 0x61, 0x74)
	if err != nil {
		return
	}
	err = en.WriteTime(z.BlinkedAt)
	if err != nil {
		err = msgp.WrapError(err, "BlinkedAt")
		return
	}
	// write "blinked"
	err = en.Append(0xa7, 0x62, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Blinked)
	if err != nil {
		err = msgp.WrapError(err, "Blinked")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Player) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 16
	// string "id"
	o = append(o, 0xde, 0x0, 0x10, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "name"
	o = append(o, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "color"
	o = append(o, 0xa5, 0x63, 0x6f, 0x6c, 0x6f, 0x72)
	o, err = z.Color.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Color")
		return
	}
	// string "joined_at"
	o = append(o, 0xa9, 0x6a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x5f, 0x61, 0x74)
	o = msgp.AppendTime(o, z.JoinedAt)
	// string "position"
	o = append(o, 0xa8, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	o, err = z.Position.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Position")
		return
	}
	// string "velocity"
	o = append(o, 0xa8, 0x76, 0x65, 0x6c, 0x6f, 0x63, 0x69, 0x74, 0x79)
	o, err = z.Velocity.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Velocity")
		return
	}
	// string "angle"
	o = append(o, 0xa5, 0x61, 0x6e, 0x67, 0x6c, 0x65)
	o = msgp.AppendFloat64(o, z.Angle)
	// string "move_dir"
	o = append(o, 0xa8, 0x6d, 0x6f, 0x76, 0x65, 0x5f, 0x64, 0x69, 0x72)
	o = msgp.AppendString(o, z.MoveDir)
	// string "turn_dir"
	o = append(o, 0xa8, 0x74, 0x75, 0x72, 0x6e, 0x5f, 0x64, 0x69, 0x72)
	o = msgp.AppendByte(o, byte(z.RotationDir))
	// string "hook"
	o = append(o, 0xa4, 0x68, 0x6f, 0x6f, 0x6b)
	if z.Hook == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Hook.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Hook")
			return
		}
	}
	// string "hooked_at"
	o = append(o, 0xa9, 0x68, 0x6f, 0x6f, 0x6b, 0x65, 0x64, 0x5f, 0x61, 0x74)
	o = msgp.AppendTime(o, z.HookedAt)
	// string "is_hooked"
	o = append(o, 0xa9, 0x69, 0x73, 0x5f, 0x68, 0x6f, 0x6f, 0x6b, 0x65, 0x64)
	o = msgp.AppendBool(o, z.IsHooked)
	// string "caught_by_id"
	o = append(o, 0xac, 0x63, 0x61, 0x75, 0x67, 0x68, 0x74, 0x5f, 0x62, 0x79, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.CaughtByID)
	// string "blinking"
	o = append(o, 0xa8, 0x62, 0x6c, 0x69, 0x6e, 0x6b, 0x69, 0x6e, 0x67)
	o = msgp.AppendBool(o, z.Blinking)
	// string "blinked_at"
	o = append(o, 0xaa, 0x62, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x5f, 0x61, 0x74)
	o = msgp.AppendTime(o, z.BlinkedAt)
	// string "blinked"
	o = append(o, 0xa7, 0x62, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64)
	o = msgp.AppendBool(o, z.Blinked)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Player) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "id":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "color":
			bts, err = z.Color.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Color")
				return
			}
		case "joined_at":
			z.JoinedAt, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "JoinedAt")
				return
			}
		case "position":
			bts, err = z.Position.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Position")
				return
			}
		case "velocity":
			bts, err = z.Velocity.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Velocity")
				return
			}
		case "angle":
			z.Angle, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Angle")
				return
			}
		case "move_dir":
			z.MoveDir, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MoveDir")
				return
			}
		case "turn_dir":
			{
				var zb0002 byte
				zb0002, bts, err = msgp.ReadByteBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "RotationDir")
					return
				}
				z.RotationDir = RotationDirection(zb0002)
			}
		case "hook":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Hook = nil
			} else {
				if z.Hook == nil {
					z.Hook = new(Hook)
				}
				bts, err = z.Hook.UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "Hook")
					return
				}
			}
		case "hooked_at":
			z.HookedAt, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "HookedAt")
				return
			}
		case "is_hooked":
			z.IsHooked, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "IsHooked")
				return
			}
		case "caught_by_id":
			z.CaughtByID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "CaughtByID")
				return
			}
		case "blinking":
			z.Blinking, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Blinking")
				return
			}
		case "blinked_at":
			z.BlinkedAt, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlinkedAt")
				return
			}
		case "blinked":
			z.Blinked, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Blinked")
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
func (z *Player) Msgsize() (s int) {
	s = 3 + 3 + msgp.StringPrefixSize + len(z.ID) + 5 + msgp.StringPrefixSize + len(z.Name) + 6 + z.Color.Msgsize() + 10 + msgp.TimeSize + 9 + z.Position.Msgsize() + 9 + z.Velocity.Msgsize() + 6 + msgp.Float64Size + 9 + msgp.StringPrefixSize + len(z.MoveDir) + 9 + msgp.ByteSize + 5
	if z.Hook == nil {
		s += msgp.NilSize
	} else {
		s += z.Hook.Msgsize()
	}
	s += 10 + msgp.TimeSize + 10 + msgp.BoolSize + 13 + msgp.StringPrefixSize + len(z.CaughtByID) + 9 + msgp.BoolSize + 11 + msgp.TimeSize + 8 + msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *RotationDirection) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 byte
		zb0001, err = dc.ReadByte()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = RotationDirection(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z RotationDirection) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteByte(byte(z))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z RotationDirection) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendByte(o, byte(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RotationDirection) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 byte
		zb0001, bts, err = msgp.ReadByteBytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = RotationDirection(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z RotationDirection) Msgsize() (s int) {
	s = msgp.ByteSize
	return
}
