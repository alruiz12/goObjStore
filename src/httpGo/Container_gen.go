package httpGo

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Container) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "objs":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Objs == nil && msz > 0 {
				z.Objs = make(map[string]Object, msz)
			} else if len(z.Objs) > 0 {
				for key, _ := range z.Objs {
					delete(z.Objs, key)
				}
			}
			for msz > 0 {
				msz--
				var xvk string
				var bzg Object
				xvk, err = dc.ReadString()
				if err != nil {
					return
				}
				err = bzg.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.Objs[xvk] = bzg
			}
		case "policy":
			z.Policy, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Container) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "name"
	err = en.Append(0x83, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "objs"
	err = en.Append(0xa4, 0x6f, 0x62, 0x6a, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.Objs)))
	if err != nil {
		return
	}
	for xvk, bzg := range z.Objs {
		err = en.WriteString(xvk)
		if err != nil {
			return
		}
		err = bzg.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "policy"
	err = en.Append(0xa6, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Policy)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Container) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "name"
	o = append(o, 0x83, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "objs"
	o = append(o, 0xa4, 0x6f, 0x62, 0x6a, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Objs)))
	for xvk, bzg := range z.Objs {
		o = msgp.AppendString(o, xvk)
		o, err = bzg.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "policy"
	o = append(o, 0xa6, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79)
	o = msgp.AppendString(o, z.Policy)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Container) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "objs":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Objs == nil && msz > 0 {
				z.Objs = make(map[string]Object, msz)
			} else if len(z.Objs) > 0 {
				for key, _ := range z.Objs {
					delete(z.Objs, key)
				}
			}
			for msz > 0 {
				var xvk string
				var bzg Object
				msz--
				xvk, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bts, err = bzg.UnmarshalMsg(bts)
				if err != nil {
					return
				}
				z.Objs[xvk] = bzg
			}
		case "policy":
			z.Policy, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *Container) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.MapHeaderSize
	if z.Objs != nil {
		for xvk, bzg := range z.Objs {
			_ = bzg
			s += msgp.StringPrefixSize + len(xvk) + bzg.Msgsize()
		}
	}
	s += 7 + msgp.StringPrefixSize + len(z.Policy)
	return
}
