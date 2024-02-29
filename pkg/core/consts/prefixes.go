package consts

type PrefixConst string

const (
	BlankPrefixConst PrefixConst = ""
	TestPrefixConst  PrefixConst = "test_"
)

func (p PrefixConst) ToString() string {
	return string(p)
}
