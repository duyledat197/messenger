package main

import (
	"github.com/dave/jennifer/jen"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"openmyth/messgener/protoc/generator"
)

const (
	extensionName = "enum"
)

func main() {
	opt := protogen.Options{}
	opt.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range p.Files {
			if !f.Generate {
				continue
			}

			gen := generator.NewGenerator(p, f, extensionName)
			var (
				enumConvertersStmts []jen.Code
			)

			for _, e := range gen.GetEnumInfos() {
				enumConvertersStmts = append(enumConvertersStmts, generateEnumConverters(e))
			}

			gen.GetJen().Add(enumConvertersStmts...)
			if len(enumConvertersStmts) > 0 {
				if err := gen.Render(); err != nil {
					grpclog.Errorf("unable to generate %s: %v", extensionName, err)
				}
			}
		}

		return nil
	})
}
