package messages

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ClientUDPMessage) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Message":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "Message")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					err = msgp.WrapError(err, "Message")
					return
				}
				switch msgp.UnsafeString(field) {
				case "type":
					{
						var zb0003 int
						zb0003, err = dc.ReadInt()
						if err != nil {
							err = msgp.WrapError(err, "Message", "T")
							return
						}
						z.Message.T = MessageType(zb0003)
					}
				case "body":
					err = z.Message.B.DecodeMsg(dc)
					if err != nil {
						err = msgp.WrapError(err, "Message", "B")
						return
					}
				default:
					err = dc.Skip()
					if err != nil {
						err = msgp.WrapError(err, "Message")
						return
					}
				}
			}
		case "id":
			z.ID, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "ID")
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
func (z *ClientUDPMessage) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Message"
	err = en.Append(0x82, 0xa7, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	if err != nil {
		return
	}
	// map header, size 2
	// write "type"
	err = en.Append(0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(int(z.Message.T))
	if err != nil {
		err = msgp.WrapError(err, "Message", "T")
		return
	}
	// write "body"
	err = en.Append(0xa4, 0x62, 0x6f, 0x64, 0x79)
	if err != nil {
		return
	}
	err = z.Message.B.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Message", "B")
		return
	}
	// write "id"
	err = en.Append(0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ClientUDPMessage) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Message"
	o = append(o, 0x82, 0xa7, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	// map header, size 2
	// string "type"
	o = append(o, 0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, int(z.Message.T))
	// string "body"
	o = append(o, 0xa4, 0x62, 0x6f, 0x64, 0x79)
	o, err = z.Message.B.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Message", "B")
		return
	}
	// string "id"
	o = append(o, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ClientUDPMessage) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Message":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Message")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					err = msgp.WrapError(err, "Message")
					return
				}
				switch msgp.UnsafeString(field) {
				case "type":
					{
						var zb0003 int
						zb0003, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							err = msgp.WrapError(err, "Message", "T")
							return
						}
						z.Message.T = MessageType(zb0003)
					}
				case "body":
					bts, err = z.Message.B.UnmarshalMsg(bts)
					if err != nil {
						err = msgp.WrapError(err, "Message", "B")
						return
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						err = msgp.WrapError(err, "Message")
						return
					}
				}
			}
		case "id":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ID")
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
func (z *ClientUDPMessage) Msgsize() (s int) {
	s = 1 + 8 + 1 + 5 + msgp.IntSize + 5 + z.Message.B.Msgsize() + 3 + msgp.StringPrefixSize + len(z.ID)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Empty) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z Empty) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 0
	_ = z
	err = en.Append(0x80)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Empty) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 0
	_ = z
	o = append(o, 0x80)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Empty) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z Empty) Msgsize() (s int) {
	s = 1
	return
}

// DecodeMsg implements msgp.Decodable
func (z *JoinGameMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Name")
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
func (z JoinGameMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "name"
	err = en.Append(0x81, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z JoinGameMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "name"
	o = append(o, 0x81, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *JoinGameMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Name")
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
func (z JoinGameMsg) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Message) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "type":
			{
				var zb0002 int
				zb0002, err = dc.ReadInt()
				if err != nil {
					err = msgp.WrapError(err, "T")
					return
				}
				z.T = MessageType(zb0002)
			}
		case "body":
			err = z.B.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "B")
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
func (z *Message) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "type"
	err = en.Append(0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(int(z.T))
	if err != nil {
		err = msgp.WrapError(err, "T")
		return
	}
	// write "body"
	err = en.Append(0xa4, 0x62, 0x6f, 0x64, 0x79)
	if err != nil {
		return
	}
	err = z.B.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "B")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Message) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "type"
	o = append(o, 0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, int(z.T))
	// string "body"
	o = append(o, 0xa4, 0x62, 0x6f, 0x64, 0x79)
	o, err = z.B.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "B")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Message) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "type":
			{
				var zb0002 int
				zb0002, bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "T")
					return
				}
				z.T = MessageType(zb0002)
			}
		case "body":
			bts, err = z.B.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "B")
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
func (z *Message) Msgsize() (s int) {
	s = 1 + 5 + msgp.IntSize + 5 + z.B.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *MessageType) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 int
		zb0001, err = dc.ReadInt()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = MessageType(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z MessageType) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteInt(int(z))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z MessageType) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt(o, int(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MessageType) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 int
		zb0001, bts, err = msgp.ReadIntBytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = MessageType(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z MessageType) Msgsize() (s int) {
	s = msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *MoveMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "dir":
			z.Dir, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Dir")
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
func (z MoveMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "dir"
	err = en.Append(0x81, 0xa3, 0x64, 0x69, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Dir)
	if err != nil {
		err = msgp.WrapError(err, "Dir")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z MoveMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "dir"
	o = append(o, 0x81, 0xa3, 0x64, 0x69, 0x72)
	o = msgp.AppendString(o, z.Dir)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MoveMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "dir":
			z.Dir, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Dir")
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
func (z MoveMsg) Msgsize() (s int) {
	s = 1 + 4 + msgp.StringPrefixSize + len(z.Dir)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PlayerBrakedMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z PlayerBrakedMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "id"
	err = en.Append(0x81, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z PlayerBrakedMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "id"
	o = append(o, 0x81, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PlayerBrakedMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z PlayerBrakedMsg) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PlayerMovedMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "dir":
			z.Dir, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Dir")
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
func (z PlayerMovedMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "id"
	err = en.Append(0x82, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "dir"
	err = en.Append(0xa3, 0x64, 0x69, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Dir)
	if err != nil {
		err = msgp.WrapError(err, "Dir")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z PlayerMovedMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "id"
	o = append(o, 0x82, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "dir"
	o = append(o, 0xa3, 0x64, 0x69, 0x72)
	o = msgp.AppendString(o, z.Dir)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PlayerMovedMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "dir":
			z.Dir, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Dir")
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
func (z PlayerMovedMsg) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 4 + msgp.StringPrefixSize + len(z.Dir)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PlayerTeleportedMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z PlayerTeleportedMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "id"
	err = en.Append(0x81, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z PlayerTeleportedMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "id"
	o = append(o, 0x81, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PlayerTeleportedMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z PlayerTeleportedMsg) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *UdpMoveMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ClientUDPMessage":
			err = z.ClientUDPMessage.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "ClientUDPMessage")
				return
			}
		case "MoveMsg":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "MoveMsg")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					err = msgp.WrapError(err, "MoveMsg")
					return
				}
				switch msgp.UnsafeString(field) {
				case "dir":
					z.MoveMsg.Dir, err = dc.ReadString()
					if err != nil {
						err = msgp.WrapError(err, "MoveMsg", "Dir")
						return
					}
				default:
					err = dc.Skip()
					if err != nil {
						err = msgp.WrapError(err, "MoveMsg")
						return
					}
				}
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
func (z *UdpMoveMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ClientUDPMessage"
	err = en.Append(0x82, 0xb0, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x55, 0x44, 0x50, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	if err != nil {
		return
	}
	err = z.ClientUDPMessage.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "ClientUDPMessage")
		return
	}
	// write "MoveMsg"
	err = en.Append(0xa7, 0x4d, 0x6f, 0x76, 0x65, 0x4d, 0x73, 0x67)
	if err != nil {
		return
	}
	// map header, size 1
	// write "dir"
	err = en.Append(0x81, 0xa3, 0x64, 0x69, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.MoveMsg.Dir)
	if err != nil {
		err = msgp.WrapError(err, "MoveMsg", "Dir")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *UdpMoveMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ClientUDPMessage"
	o = append(o, 0x82, 0xb0, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x55, 0x44, 0x50, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	o, err = z.ClientUDPMessage.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ClientUDPMessage")
		return
	}
	// string "MoveMsg"
	o = append(o, 0xa7, 0x4d, 0x6f, 0x76, 0x65, 0x4d, 0x73, 0x67)
	// map header, size 1
	// string "dir"
	o = append(o, 0x81, 0xa3, 0x64, 0x69, 0x72)
	o = msgp.AppendString(o, z.MoveMsg.Dir)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *UdpMoveMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ClientUDPMessage":
			bts, err = z.ClientUDPMessage.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "ClientUDPMessage")
				return
			}
		case "MoveMsg":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MoveMsg")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					err = msgp.WrapError(err, "MoveMsg")
					return
				}
				switch msgp.UnsafeString(field) {
				case "dir":
					z.MoveMsg.Dir, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "MoveMsg", "Dir")
						return
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						err = msgp.WrapError(err, "MoveMsg")
						return
					}
				}
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
func (z *UdpMoveMsg) Msgsize() (s int) {
	s = 1 + 17 + z.ClientUDPMessage.Msgsize() + 8 + 1 + 4 + msgp.StringPrefixSize + len(z.MoveMsg.Dir)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *YourIDMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z YourIDMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "id"
	err = en.Append(0x81, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z YourIDMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "id"
	o = append(o, 0x81, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *YourIDMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z YourIDMsg) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID)
	return
}