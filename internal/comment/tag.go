package comment

type Tag string

const tagSep = ":"

const (
	NoTag    Tag = ""
	BoolTag  Tag = "bool"
	StrTag   Tag = "str"
	IntTag   Tag = "int"
	FloatTag Tag = "float"
	SeqTag   Tag = "seq"
	MapTag   Tag = "map"
)

func Tags() []Tag {
	return []Tag{
		BoolTag,
		StrTag,
		IntTag,
		FloatTag,
		SeqTag,
		MapTag,
	}
}
