package main

import (
	"github.com/dave/jennifer/jen"

	pb "openmyth/messgener/pb/base"
)

// getEncodeLib returns the library corresponding to the given event encoding type.
func getEncodeLib(event pb.EventEncodeType) string {
	switch event {
	case pb.EventEncodeType_EventEncodeType_JSON:
		return "encoding/json"
	case pb.EventEncodeType_EventEncodeType_Protobuf:
		return "google.golang.org/protobuf/proto"
	}

	return ""
}

// generateTopicConst generates a constant statement for the given event topic.
func generateTopicConst(event *pb.EventOptions) *jen.Statement {
	return jen.Id(event.Topic).Op("=").Lit(event.Topic).Comment("event")
}

// generateMarshalFn generates a function comment for the given message name and event options.
func generateMarshalFn(msgName string, event *pb.EventOptions) *jen.Statement {
	return jen.
		Comment("Marshal marshals a "+msgName+" to bytes").Line().
		Func().Params(
		jen.Id("x").Op("*").Id(msgName),
	).
		Id("Marshal").
		Params().
		Params(jen.Index().Byte(), jen.Error()).
		Block(
			jen.Id("b").Op(",").Err().Op(":=").Qual(getEncodeLib(event.EncodeType), "Marshal").Call(jen.Id("x")),
			jen.If(jen.Err().Op("!=").Nil()).
				Block(
					jen.Return(jen.Nil(), jen.Err()),
				).Line().
				Return(jen.Id("b"), jen.Nil()),
		).
		Line()
}

// generateUnmarshalFn generates a function comment for the given message name and event options.
func generateUnmarshalFn(msgName string, event *pb.EventOptions) *jen.Statement {
	return jen.
		Comment("Unmarshal unmarshals from bytes to a " + msgName).Line().
		Func().Params(
		jen.Id("x").Id(msgName),
	).
		Id("Unmarshal").
		Params(jen.Id("b").Index().Byte()).
		Params(jen.Error()).
		Block(
			jen.If(
				jen.Err().Op(":=").Qual(getEncodeLib(event.EncodeType), "Unmarshal").Call(jen.Id("b"), jen.Id("&x")).Op(";").
					Err().Op("!=").Nil()).
				Block(
					jen.Return(jen.Err()),
				).Line().
				Return(jen.Nil()),
		).
		Line()
}
