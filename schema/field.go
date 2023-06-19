package schema

const (
	Bool   string = "BOOL"
	Int    string = "INT"
	Uint   string = "INT UNSIGNED"
	Float  string = "FLOAT"
	String string = "VARCHAR(255)"
	//Time   string = "time"
	//Bytes  string = "bytes"
)

type Field struct {
	Name          string
	NameDB        string
	Type          string
	PrimaryKey    bool
	ScenarioField bool
}

var GoTypesToDbTypes = map[string]string{
	"int":     Int,
	"int8":    Int,
	"int16":   Int,
	"int32":   Int,
	"int64":   Int,
	"uint":    Uint,
	"uint8":   Uint,
	"uint16":  Uint,
	"uint32":  Uint,
	"uint64":  Uint,
	"float32": Float,
	"float64": Float,
	"string":  String,
	"bool":    Bool,
}
