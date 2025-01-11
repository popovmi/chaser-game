package game

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"time"

	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Portal) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "link_id":
			z.LinkID, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "LinkID")
				return
			}
		case "last_used_at":
			z.LastUsedAt, err = dc.ReadTime()
			if err != nil {
				err = msgp.WrapError(err, "LastUsedAt")
				return
			}
		case "pos":
			err = z.Pos.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Pos")
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
func (z *Portal) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "id"
	err = en.Append(0x84, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "link_id"
	err = en.Append(0xa7, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.LinkID)
	if err != nil {
		err = msgp.WrapError(err, "LinkID")
		return
	}
	// write "last_used_at"
	err = en.Append(0xac, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74)
	if err != nil {
		return
	}
	err = en.WriteTime(z.LastUsedAt)
	if err != nil {
		err = msgp.WrapError(err, "LastUsedAt")
		return
	}
	// write "pos"
	err = en.Append(0xa3, 0x70, 0x6f, 0x73)
	if err != nil {
		return
	}
	err = z.Pos.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Pos")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Portal) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "id"
	o = append(o, 0x84, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "link_id"
	o = append(o, 0xa7, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.LinkID)
	// string "last_used_at"
	o = append(o, 0xac, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74)
	o = msgp.AppendTime(o, z.LastUsedAt)
	// string "pos"
	o = append(o, 0xa3, 0x70, 0x6f, 0x73)
	o, err = z.Pos.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Pos")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Portal) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "link_id":
			z.LinkID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "LinkID")
				return
			}
		case "last_used_at":
			z.LastUsedAt, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "LastUsedAt")
				return
			}
		case "pos":
			bts, err = z.Pos.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Pos")
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
func (z *Portal) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 8 + msgp.StringPrefixSize + len(z.LinkID) + 13 + msgp.TimeSize + 4 + z.Pos.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PortalLink) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "portals":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "PortalIDs")
				return
			}
			if cap(z.PortalIDs) >= int(zb0002) {
				z.PortalIDs = (z.PortalIDs)[:zb0002]
			} else {
				z.PortalIDs = make([]string, zb0002)
			}
			for za0001 := range z.PortalIDs {
				z.PortalIDs[za0001], err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "PortalIDs", za0001)
					return
				}
			}
		case "lu":
			var zb0003 uint32
			zb0003, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "LastUsed")
				return
			}
			if z.LastUsed == nil {
				z.LastUsed = make(map[string]time.Time, zb0003)
			} else if len(z.LastUsed) > 0 {
				for key := range z.LastUsed {
					delete(z.LastUsed, key)
				}
			}
			for zb0003 > 0 {
				zb0003--
				var za0002 string
				var za0003 time.Time
				za0002, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "LastUsed")
					return
				}
				za0003, err = dc.ReadTime()
				if err != nil {
					err = msgp.WrapError(err, "LastUsed", za0002)
					return
				}
				z.LastUsed[za0002] = za0003
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
func (z *PortalLink) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "id"
	err = en.Append(0x83, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "portals"
	err = en.Append(0xa7, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.PortalIDs)))
	if err != nil {
		err = msgp.WrapError(err, "PortalIDs")
		return
	}
	for za0001 := range z.PortalIDs {
		err = en.WriteString(z.PortalIDs[za0001])
		if err != nil {
			err = msgp.WrapError(err, "PortalIDs", za0001)
			return
		}
	}
	// write "lu"
	err = en.Append(0xa2, 0x6c, 0x75)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.LastUsed)))
	if err != nil {
		err = msgp.WrapError(err, "LastUsed")
		return
	}
	for za0002, za0003 := range z.LastUsed {
		err = en.WriteString(za0002)
		if err != nil {
			err = msgp.WrapError(err, "LastUsed")
			return
		}
		err = en.WriteTime(za0003)
		if err != nil {
			err = msgp.WrapError(err, "LastUsed", za0002)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PortalLink) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "id"
	o = append(o, 0x83, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "portals"
	o = append(o, 0xa7, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.PortalIDs)))
	for za0001 := range z.PortalIDs {
		o = msgp.AppendString(o, z.PortalIDs[za0001])
	}
	// string "lu"
	o = append(o, 0xa2, 0x6c, 0x75)
	o = msgp.AppendMapHeader(o, uint32(len(z.LastUsed)))
	for za0002, za0003 := range z.LastUsed {
		o = msgp.AppendString(o, za0002)
		o = msgp.AppendTime(o, za0003)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PortalLink) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "portals":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PortalIDs")
				return
			}
			if cap(z.PortalIDs) >= int(zb0002) {
				z.PortalIDs = (z.PortalIDs)[:zb0002]
			} else {
				z.PortalIDs = make([]string, zb0002)
			}
			for za0001 := range z.PortalIDs {
				z.PortalIDs[za0001], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "PortalIDs", za0001)
					return
				}
			}
		case "lu":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "LastUsed")
				return
			}
			if z.LastUsed == nil {
				z.LastUsed = make(map[string]time.Time, zb0003)
			} else if len(z.LastUsed) > 0 {
				for key := range z.LastUsed {
					delete(z.LastUsed, key)
				}
			}
			for zb0003 > 0 {
				var za0002 string
				var za0003 time.Time
				zb0003--
				za0002, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "LastUsed")
					return
				}
				za0003, bts, err = msgp.ReadTimeBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "LastUsed", za0002)
					return
				}
				z.LastUsed[za0002] = za0003
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
func (z *PortalLink) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 8 + msgp.ArrayHeaderSize
	for za0001 := range z.PortalIDs {
		s += msgp.StringPrefixSize + len(z.PortalIDs[za0001])
	}
	s += 3 + msgp.MapHeaderSize
	if z.LastUsed != nil {
		for za0002, za0003 := range z.LastUsed {
			_ = za0003
			s += msgp.StringPrefixSize + len(za0002) + msgp.TimeSize
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PortalNetwork) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "portals":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "Portals")
				return
			}
			if z.Portals == nil {
				z.Portals = make(map[string]*Portal, zb0002)
			} else if len(z.Portals) > 0 {
				for key := range z.Portals {
					delete(z.Portals, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 *Portal
				za0001, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Portals")
					return
				}
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						err = msgp.WrapError(err, "Portals", za0001)
						return
					}
					za0002 = nil
				} else {
					if za0002 == nil {
						za0002 = new(Portal)
					}
					err = za0002.DecodeMsg(dc)
					if err != nil {
						err = msgp.WrapError(err, "Portals", za0001)
						return
					}
				}
				z.Portals[za0001] = za0002
			}
		case "portal_links":
			var zb0003 uint32
			zb0003, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "Links")
				return
			}
			if z.Links == nil {
				z.Links = make(map[string]*PortalLink, zb0003)
			} else if len(z.Links) > 0 {
				for key := range z.Links {
					delete(z.Links, key)
				}
			}
			for zb0003 > 0 {
				zb0003--
				var za0003 string
				var za0004 *PortalLink
				za0003, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Links")
					return
				}
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						err = msgp.WrapError(err, "Links", za0003)
						return
					}
					za0004 = nil
				} else {
					if za0004 == nil {
						za0004 = new(PortalLink)
					}
					err = za0004.DecodeMsg(dc)
					if err != nil {
						err = msgp.WrapError(err, "Links", za0003)
						return
					}
				}
				z.Links[za0003] = za0004
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
func (z *PortalNetwork) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "portals"
	err = en.Append(0x82, 0xa7, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Portals)))
	if err != nil {
		err = msgp.WrapError(err, "Portals")
		return
	}
	for za0001, za0002 := range z.Portals {
		err = en.WriteString(za0001)
		if err != nil {
			err = msgp.WrapError(err, "Portals")
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
				err = msgp.WrapError(err, "Portals", za0001)
				return
			}
		}
	}
	// write "portal_links"
	err = en.Append(0xac, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Links)))
	if err != nil {
		err = msgp.WrapError(err, "Links")
		return
	}
	for za0003, za0004 := range z.Links {
		err = en.WriteString(za0003)
		if err != nil {
			err = msgp.WrapError(err, "Links")
			return
		}
		if za0004 == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = za0004.EncodeMsg(en)
			if err != nil {
				err = msgp.WrapError(err, "Links", za0003)
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PortalNetwork) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "portals"
	o = append(o, 0x82, 0xa7, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Portals)))
	for za0001, za0002 := range z.Portals {
		o = msgp.AppendString(o, za0001)
		if za0002 == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = za0002.MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "Portals", za0001)
				return
			}
		}
	}
	// string "portal_links"
	o = append(o, 0xac, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Links)))
	for za0003, za0004 := range z.Links {
		o = msgp.AppendString(o, za0003)
		if za0004 == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = za0004.MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "Links", za0003)
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PortalNetwork) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "portals":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Portals")
				return
			}
			if z.Portals == nil {
				z.Portals = make(map[string]*Portal, zb0002)
			} else if len(z.Portals) > 0 {
				for key := range z.Portals {
					delete(z.Portals, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 *Portal
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Portals")
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
						za0002 = new(Portal)
					}
					bts, err = za0002.UnmarshalMsg(bts)
					if err != nil {
						err = msgp.WrapError(err, "Portals", za0001)
						return
					}
				}
				z.Portals[za0001] = za0002
			}
		case "portal_links":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Links")
				return
			}
			if z.Links == nil {
				z.Links = make(map[string]*PortalLink, zb0003)
			} else if len(z.Links) > 0 {
				for key := range z.Links {
					delete(z.Links, key)
				}
			}
			for zb0003 > 0 {
				var za0003 string
				var za0004 *PortalLink
				zb0003--
				za0003, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Links")
					return
				}
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					za0004 = nil
				} else {
					if za0004 == nil {
						za0004 = new(PortalLink)
					}
					bts, err = za0004.UnmarshalMsg(bts)
					if err != nil {
						err = msgp.WrapError(err, "Links", za0003)
						return
					}
				}
				z.Links[za0003] = za0004
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
func (z *PortalNetwork) Msgsize() (s int) {
	s = 1 + 8 + msgp.MapHeaderSize
	if z.Portals != nil {
		for za0001, za0002 := range z.Portals {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001)
			if za0002 == nil {
				s += msgp.NilSize
			} else {
				s += za0002.Msgsize()
			}
		}
	}
	s += 13 + msgp.MapHeaderSize
	if z.Links != nil {
		for za0003, za0004 := range z.Links {
			_ = za0004
			s += msgp.StringPrefixSize + len(za0003)
			if za0004 == nil {
				s += msgp.NilSize
			} else {
				s += za0004.Msgsize()
			}
		}
	}
	return
}
