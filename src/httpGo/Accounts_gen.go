package httpGo

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Accounts) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "accounts":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Accounts == nil && msz > 0 {
				z.Accounts = make(map[string]Account, msz)
			} else if len(z.Accounts) > 0 {
				for key, _ := range z.Accounts {
					delete(z.Accounts, key)
				}
			}
			for msz > 0 {
				msz--
				var xvk string
				var bzg Account
				xvk, err = dc.ReadString()
				if err != nil {
					return
				}
				err = bzg.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.Accounts[xvk] = bzg
			}
		case "num":
			z.Num, err = dc.ReadInt()
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
func (z *Accounts) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "accounts"
	err = en.Append(0x82, 0xa8, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.Accounts)))
	if err != nil {
		return
	}
	for xvk, bzg := range z.Accounts {
		err = en.WriteString(xvk)
		if err != nil {
			return
		}
		err = bzg.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "num"
	err = en.Append(0xa3, 0x6e, 0x75, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Num)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Accounts) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "accounts"
	o = append(o, 0x82, 0xa8, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Accounts)))
	for xvk, bzg := range z.Accounts {
		o = msgp.AppendString(o, xvk)
		o, err = bzg.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "num"
	o = append(o, 0xa3, 0x6e, 0x75, 0x6d)
	o = msgp.AppendInt(o, z.Num)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Accounts) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "accounts":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Accounts == nil && msz > 0 {
				z.Accounts = make(map[string]Account, msz)
			} else if len(z.Accounts) > 0 {
				for key, _ := range z.Accounts {
					delete(z.Accounts, key)
				}
			}
			for msz > 0 {
				var xvk string
				var bzg Account
				msz--
				xvk, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bts, err = bzg.UnmarshalMsg(bts)
				if err != nil {
					return
				}
				z.Accounts[xvk] = bzg
			}
		case "num":
			z.Num, bts, err = msgp.ReadIntBytes(bts)
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

func (z *Accounts) Msgsize() (s int) {
	s = 1 + 9 + msgp.MapHeaderSize
	if z.Accounts != nil {
		for xvk, bzg := range z.Accounts {
			_ = bzg
			s += msgp.StringPrefixSize + len(xvk) + bzg.Msgsize()
		}
	}
	s += 4 + msgp.IntSize
	return
}
