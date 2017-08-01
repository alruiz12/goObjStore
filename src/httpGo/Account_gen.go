package httpGo

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Account) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "containers":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Containers == nil && msz > 0 {
				z.Containers = make(map[string]Container, msz)
			} else if len(z.Containers) > 0 {
				for key, _ := range z.Containers {
					delete(z.Containers, key)
				}
			}
			for msz > 0 {
				msz--
				var xvk string
				var bzg Container
				xvk, err = dc.ReadString()
				if err != nil {
					return
				}
				err = bzg.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.Containers[xvk] = bzg
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
func (z *Account) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "name"
	err = en.Append(0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "containers"
	err = en.Append(0xaa, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.Containers)))
	if err != nil {
		return
	}
	for xvk, bzg := range z.Containers {
		err = en.WriteString(xvk)
		if err != nil {
			return
		}
		err = bzg.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Account) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "name"
	o = append(o, 0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "containers"
	o = append(o, 0xaa, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Containers)))
	for xvk, bzg := range z.Containers {
		o = msgp.AppendString(o, xvk)
		o, err = bzg.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Account) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "containers":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Containers == nil && msz > 0 {
				z.Containers = make(map[string]Container, msz)
			} else if len(z.Containers) > 0 {
				for key, _ := range z.Containers {
					delete(z.Containers, key)
				}
			}
			for msz > 0 {
				var xvk string
				var bzg Container
				msz--
				xvk, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bts, err = bzg.UnmarshalMsg(bts)
				if err != nil {
					return
				}
				z.Containers[xvk] = bzg
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

func (z *Account) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 11 + msgp.MapHeaderSize
	if z.Containers != nil {
		for xvk, bzg := range z.Containers {
			_ = bzg
			s += msgp.StringPrefixSize + len(xvk) + bzg.Msgsize()
		}
	}
	return
}
