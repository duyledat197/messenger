package main

import (
	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
)

// generateEnumConverters generates a function comment for the given enum.
// It takes a pointer to a protogen.Enum as a parameter and returns a *jen.Statement.
func generateEnumConverters(e *protogen.Enum) *jen.Statement {
	enumName := string(e.Desc.Name())
	return jen.Func().Params(jen.Id("x").Id(enumName)).Id("FromString").Params(
		jen.Id("str").String(),
	).Params(jen.Id(enumName), jen.Error()).Block(
		jen.Switch(jen.Id("str")).BlockFunc(
			func(g *jen.Group) {
				for _, val := range e.Values {
					name := val.GoIdent.GoName
					g.Case(jen.Id(name).Dot("String").Call()).Block(
						jen.Return(jen.Id(name), jen.Nil()),
					)
				}

				g.Default().
					Block(
						jen.Return(
							jen.Id(enumName).Call(jen.Lit(0)),
							jen.Qual("fmt", "Errorf").
								Call(jen.Lit("unknown %s"), jen.Id("str")),
						),
					)
			},
		),
	)
}
