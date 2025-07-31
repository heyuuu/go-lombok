package astkit

import (
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"strings"
)

var dumpCfg = printer.Config{
	Mode:     printer.UseSpaces | printer.TabIndent,
	Tabwidth: 8,
}

func PrintNode(node ast.Node) string {
	var buf strings.Builder
	err := dumpCfg.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
