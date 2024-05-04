package generator

import (
	"strings"

	"github.com/modern-go/reflect2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

// getOption retrieves an option from a protocol buffer message.
// It takes in a message of type protoreflect.Message and an extension of type *protoimpl.ExtensionInfo.
// It returns a proto.Message and a bool indicating whether the option was found or not.
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

func getGoPackageName(file *descriptorpb.FileDescriptorProto) (impPath string, pkg string, ok bool) {
	opt := file.GetOptions().GetGoPackage()
	if opt == "" {
		impPath = file.GetPackage()
		pkg = impPath
		return impPath, pkg, false
	}
	// A semicolon-delimited suffix delimits the import path and package name.
	sc := strings.Index(opt, ";")
	if sc >= 0 {
		return opt[:sc], opt[sc+1:], true
	}
	// The presence of a slash implies there's an import path.
	slash := strings.LastIndex(opt, "/")
	if slash >= 0 {
		return opt, opt[slash+1:], true
	}
	return "", opt, true
}
