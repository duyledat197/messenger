package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	opt := protogen.Options{}

	opt.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range p.Files {
			if !f.Generate {
				continue
			}

			// g := generator{f, p}
			// g.GenerateMethods()
			// g.GenerateEnums()
		}
		return nil
	})
}
