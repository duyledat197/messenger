package main

import (
	"github.com/dave/jennifer/jen"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	pb "openmyth/messgener/pb/base"
	"openmyth/messgener/protoc/generator"
)

const (
	extensionName = "event"
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

			msgEvents, err := gen.GetOptionInfos(pb.E_Event)
			if err != nil {
				grpclog.Errorf("unable to generate event: %v", err)
			}

			var (
				topicConstStmts []jen.Code
				marshalStmts    []jen.Code
				unmarshalStmts  []jen.Code
			)
			for msgName, msgEvent := range msgEvents {
				event, ok := msgEvent.(*pb.EventOptions)
				if !ok {
					continue
				}
				topicConstStmts = append(topicConstStmts, generateTopicConst(event))
				marshalStmts = append(marshalStmts, generateMarshalFn(msgName, event))
				unmarshalStmts = append(unmarshalStmts, generateUnmarshalFn(msgName, event))
			}

			gen.GetJen().Const().Defs(topicConstStmts...)

			gen.GetJen().Add(marshalStmts...)
			gen.GetJen().Add(unmarshalStmts...)

			if err := gen.Render(); err != nil {
				grpclog.Errorf("unable to generate event: %v", err)
			}
		}

		return nil
	})
}
