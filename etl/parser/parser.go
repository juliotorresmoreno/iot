package parser

type provider interface {
	Encode(interface{}) (string, error)
	Decode(string, interface{}) error
}

type Parser struct {
	provider provider
}

type ParserFormat int64

var (
	BSON ParserFormat = 0
	JSON ParserFormat = 1
)

func MakeParser(format ...ParserFormat) *Parser {
	p := &Parser{}
	if len(format) == 0 || format[0] == BSON {
		p.provider = NewBSONProvider()
	}
	return p
}

func (el *Parser) Encode(value interface{}) (string, error) {
	return el.provider.Encode(value)
}

func (el *Parser) Decode(value string, obj interface{}) error {
	return el.provider.Decode(value, obj)
}
