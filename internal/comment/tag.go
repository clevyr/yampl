package comment

type Tag string

func (t Tag) ToYaml() string {
	if t == DynamicTag {
		return ""
	}
	return "!!" + string(t)
}

const tagSep = ":"

var (
	DynamicTag Tag = ""
	BoolTag    Tag = "bool"
	StrTag     Tag = "str"
	IntTag     Tag = "int"
	FloatTag   Tag = "float"
	SeqTag     Tag = "seq"
	MapTag     Tag = "map"
)

var tags = []Tag{
	BoolTag,
	StrTag,
	IntTag,
	FloatTag,
	SeqTag,
	MapTag,
}
