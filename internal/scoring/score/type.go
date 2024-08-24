package score

// -----------------------------------------------------------------------------

//go:generate go run github.com/dmarkham/enumer -type=Type
type Type uint8

const (
	Default Type = iota
	SubScore
	ByData
	Original
	ByTest
)

var colNames = map[Type]string{
	Default:  "Score",
	SubScore: "SubScore",
	ByData:   "by Data",
	ByTest:   "by Test",
	Original: "Original",
}

func (t Type) ColHeader() string {
	return colNames[t]
}
