package main

import (
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"openmyth/messgener/protoc/generator"
)

const (
	extensionName = "api-info"
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

			isGen := false
			for svcName, methods := range gen.GetMethodDescInfos() {
				stmt := generateConstAPIByService(svcName, methods)
				if stmt != nil {
					gen.GetJen().Add(stmt)
					isGen = true
				}
			}
			if isGen {
				if err := gen.Render(); err != nil {
					grpclog.Errorf("unable to generate %s: %v", extensionName, err)
				}
			}
		}

		return nil
	})

}
