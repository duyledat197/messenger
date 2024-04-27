package main

import (
	"github.com/modern-go/reflect2"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type generator struct {
	plugin *protogen.Plugin
	file   *protogen.File
}

func (g *generator) GenerateEvent() {
	for _, msg := range g.file.Proto.MessageType {

		// TODO: generate base protos for using extensions
		_, _ = getOption(msg.ProtoReflect(), &protoimpl.ExtensionInfo{})

	}

}

func getOption(message protoreflect.Message, extension *protoimpl.ExtensionInfo) (proto.Message, bool) {
	ref := message.Get(extension.TypeDescriptor())
	if !ref.IsValid() {
		return nil, false
	}
	inter := ref.Message().Interface()
	if inter == nil {
		return nil, false
	}
	if reflect2.IsNil(inter) {
		return nil, false
	}
	return inter, true
}
