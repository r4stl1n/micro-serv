package consts

type VersionConst string

const (
	VersionV1Const VersionConst = "v1"
)

func (v VersionConst) ToString() string {
	return string(v)
}
