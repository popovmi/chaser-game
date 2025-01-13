package game

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Game) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "players":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "Players")
				return
			}
			if z.Players == nil {
				z.Players = make(map[string]*Player, zb0002)
			} else if len(z.Players) > 0 {
				for key := range z.Players {
					delete(z.Players, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 *Player
				za0001, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Players")
					return
				}
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						err = msgp.WrapError(err, "Players", za0001)
						return
					}
					za0002 = nil
				} else {
					if za0002 == nil {
						za0002 = new(Player)
					}
					err = za0002.DecodeMsg(dc)
					if err != nil {
						err = msgp.WrapError(err, "Players", za0001)
						return
					}
				}
				z.Players[za0001] = za0002
			}
		case "portal_network":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "PortalNetwork")
					return
				}
				z.PortalNetwork = nil
			} else {
				if z.PortalNetwork == nil {
					z.PortalNetwork = new(PortalNetwork)
				}
				err = z.PortalNetwork.DecodeMsg(dc)
				if err != nil {
					err = msgp.WrapError(err, "PortalNetwork")
					return
				}
			}
		case "bricks":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "Bricks")
					return
				}
				z.Bricks = nil
			} else {
				var zb0003 uint32
				zb0003, err = dc.ReadArrayHeader()
				if err != nil {
					err = msgp.WrapError(err, "Bricks")
					return
				}
				if z.Bricks != nil && cap(z.Bricks) >= int(zb0003) {
					z.Bricks = (z.Bricks)[:zb0003]
				} else {
					z.Bricks = make([]*Brick, zb0003)
				}
				for za0003 := range z.Bricks {
					if dc.IsNil() {
						err = dc.ReadNil()
						if err != nil {
							err = msgp.WrapError(err, "Bricks", za0003)
							return
						}
						z.Bricks[za0003] = nil
					} else {
						if z.Bricks[za0003] == nil {
							z.Bricks[za0003] = new(Brick)
						}
						err = z.Bricks[za0003].DecodeMsg(dc)
						if err != nil {
							err = msgp.WrapError(err, "Bricks", za0003)
							return
						}
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
func (z *Game) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "players"
	err = en.Append(0x83, 0xa7, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Players)))
	if err != nil {
		err = msgp.WrapError(err, "Players")
		return
	}
	for za0001, za0002 := range z.Players {
		err = en.WriteString(za0001)
		if err != nil {
			err = msgp.WrapError(err, "Players")
			return
		}
		if za0002 == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = za0002.EncodeMsg(en)
			if err != nil {
				err = msgp.WrapError(err, "Players", za0001)
				return
			}
		}
	}
	// write "portal_network"
	err = en.Append(0xae, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b)
	if err != nil {
		return
	}
	if z.PortalNetwork == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.PortalNetwork.EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "PortalNetwork")
			return
		}
	}
	// write "bricks"
	err = en.Append(0xa6, 0x62, 0x72, 0x69, 0x63, 0x6b, 0x73)
	if err != nil {
		return
	}
	if z.Bricks == nil { // allownil: if nil
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteArrayHeader(uint32(len(z.Bricks)))
		if err != nil {
			err = msgp.WrapError(err, "Bricks")
			return
		}
		for za0003 := range z.Bricks {
			if z.Bricks[za0003] == nil {
				err = en.WriteNil()
				if err != nil {
					return
				}
			} else {
				err = z.Bricks[za0003].EncodeMsg(en)
				if err != nil {
					err = msgp.WrapError(err, "Bricks", za0003)
					return
				}
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Game) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "players"
	o = append(o, 0x83, 0xa7, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Players)))
	for za0001, za0002 := range z.Players {
		o = msgp.AppendString(o, za0001)
		if za0002 == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = za0002.MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "Players", za0001)
				return
			}
		}
	}
	// string "portal_network"
	o = append(o, 0xae, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b)
	if z.PortalNetwork == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.PortalNetwork.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "PortalNetwork")
			return
		}
	}
	// string "bricks"
	o = append(o, 0xa6, 0x62, 0x72, 0x69, 0x63, 0x6b, 0x73)
	if z.Bricks == nil { // allownil: if nil
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Bricks)))
		for za0003 := range z.Bricks {
			if z.Bricks[za0003] == nil {
				o = msgp.AppendNil(o)
			} else {
				o, err = z.Bricks[za0003].MarshalMsg(o)
				if err != nil {
					err = msgp.WrapError(err, "Bricks", za0003)
					return
				}
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Game) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "players":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Players")
				return
			}
			if z.Players == nil {
				z.Players = make(map[string]*Player, zb0002)
			} else if len(z.Players) > 0 {
				for key := range z.Players {
					delete(z.Players, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 *Player
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Players")
					return
				}
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					za0002 = nil
				} else {
					if za0002 == nil {
						za0002 = new(Player)
					}
					bts, err = za0002.UnmarshalMsg(bts)
					if err != nil {
						err = msgp.WrapError(err, "Players", za0001)
						return
					}
				}
				z.Players[za0001] = za0002
			}
		case "portal_network":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.PortalNetwork = nil
			} else {
				if z.PortalNetwork == nil {
					z.PortalNetwork = new(PortalNetwork)
				}
				bts, err = z.PortalNetwork.UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "PortalNetwork")
					return
				}
			}
		case "bricks":
			if msgp.IsNil(bts) {
				bts = bts[1:]
				z.Bricks = nil
			} else {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Bricks")
					return
				}
				if z.Bricks != nil && cap(z.Bricks) >= int(zb0003) {
					z.Bricks = (z.Bricks)[:zb0003]
				} else {
					z.Bricks = make([]*Brick, zb0003)
				}
				for za0003 := range z.Bricks {
					if msgp.IsNil(bts) {
						bts, err = msgp.ReadNilBytes(bts)
						if err != nil {
							return
						}
						z.Bricks[za0003] = nil
					} else {
						if z.Bricks[za0003] == nil {
							z.Bricks[za0003] = new(Brick)
						}
						bts, err = z.Bricks[za0003].UnmarshalMsg(bts)
						if err != nil {
							err = msgp.WrapError(err, "Bricks", za0003)
							return
						}
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
func (z *Game) Msgsize() (s int) {
	s = 1 + 8 + msgp.MapHeaderSize
	if z.Players != nil {
		for za0001, za0002 := range z.Players {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001)
			if za0002 == nil {
				s += msgp.NilSize
			} else {
				s += za0002.Msgsize()
			}
		}
	}
	s += 15
	if z.PortalNetwork == nil {
		s += msgp.NilSize
	} else {
		s += z.PortalNetwork.Msgsize()
	}
	s += 7 + msgp.ArrayHeaderSize
	for za0003 := range z.Bricks {
		if z.Bricks[za0003] == nil {
			s += msgp.NilSize
		} else {
			s += z.Bricks[za0003].Msgsize()
		}
	}
	return
}
