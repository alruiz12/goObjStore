package httpGo
//go:generate msgp
type Accounts struct{
	Accounts map[string]Account	`msg:"accounts"`
	Num int		`msg:"num"`
}