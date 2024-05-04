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

			for msgName, msgEvent := range msgEvents {
				event, ok := msgEvent.(*pb.EventOptions)
				if !ok {
					continue
				}
				generateEvents(gen.GetJen(), msgName, event)
			}

			if err := gen.Render(); err != nil {
				grpclog.Errorf("unable to generate event: %v", err)
			}
		}

		return nil
	})
}

func generateEvents(gen *jen.File, msgName string, event *pb.EventOptions) {
	gen.Const().Id(event.Topic).Op("=").Lit(event.Topic).Comment("event").Line()
	var (
		lib = ""
	)
	switch event.EncodeType {
	case pb.EventEncodeType_EventEncodeType_JSON:
		lib = "encoding/json"
	case pb.EventEncodeType_EventEncodeType_Protobuf:
		lib = "google.golang.org/protobuf/proto"
	}

	// generate marshal method
	gen.
		Comment("Marshal marshals a "+msgName+" to bytes").Line().
		Func().Params(
		jen.Id("x").Op("*").Id(msgName),
	).
		Id("Marshal").
		Params().
		Params(jen.Index().Byte(), jen.Error()).
		Block(
			jen.Id("b").Op(",").Err().Op(":=").Qual(lib, "Marshal").Call(jen.Id("x")),
			jen.If(jen.Err().Op("!=").Nil()).
				Block(
					jen.Return(jen.Nil(), jen.Err()),
				).Line().
				Return(jen.Id("b"), jen.Nil()),
		).
		Line()

		// generate unmarshal method
	gen.
		Comment("Unmarshal unmarshals from bytes to a " + msgName).Line().
		Func().Params(
		jen.Id("x").Id(msgName),
	).
		Id("Unmarshal").
		Params(jen.Id("b").Index().Byte()).
		Params(jen.Error()).
		Block(
			jen.If(
				jen.Err().Op(":=").Qual(lib, "Unmarshal").Call(jen.Id("b"), jen.Id("&x")).Op(";").
					Err().Op("!=").Nil()).
				Block(
					jen.Return(jen.Err()),
				).Line().
				Return(jen.Nil()),
		).
		Line()
}
