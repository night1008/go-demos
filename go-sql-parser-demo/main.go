package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/akito0107/xsqlparser"
	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/k0kubun/pp"
)

type ClickhouseGenericSQLDialect struct {
}

func (*ClickhouseGenericSQLDialect) IsIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '@'
}

func (*ClickhouseGenericSQLDialect) IsIdentifierPart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '@' || r == '_'
}

func (*ClickhouseGenericSQLDialect) IsDelimitedIdentifierStart(r rune) bool {
	return r == '`'
}

func main() {
	// 必须指定分号
	// str := "SELECT toDate(aa.`#vp@aaa`) + cast(bbb as int);"
	// str := "SELECT `aa` + 1;"
	// str := "SELECT a, case `aa` when 1 then 'a' else 'b' end;"
	str := "SELECT todate(`aa`, 'aa');"
	parser, err := xsqlparser.NewParser(bytes.NewBufferString(str), &ClickhouseGenericSQLDialect{})
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := parser.ParseStatement()
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(stmt)

	var list []sqlast.Node
	sqlast.Inspect(stmt, func(node sqlast.Node) bool {
		switch node.(type) {
		case nil:
			return false
		case sqlast.SQLSelectItem:
			fmt.Println("===> aaa", node.ToSQLString())
			return false
		default:
			list = append(list, node)
			return true
		}
	})
	// pp.Println(list)
}
