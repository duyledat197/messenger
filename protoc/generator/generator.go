package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Generator struct {
	plugin *protogen.Plugin
	file   *protogen.File

	jen *jen.File

	version       string
	extensionName string
}

// NewGenerator creates a new Generator instance with the specified plugin, file, and extensionName.
func NewGenerator(plugin *protogen.Plugin, file *protogen.File, extensionName string) *Generator {
	_, packageName, _ := getGoPackageName(file.Proto)
	gf := jen.NewFile(packageName)

	compilerVersion := "unknown"
	if v := plugin.Request.GetCompilerVersion(); v != nil {
		compilerVersion = fmt.Sprintf("%d.%d.%d", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}

	gf.HeaderComment(fmt.Sprintf("Code generated by protoc-gen-%v. DO NOT EDIT.", extensionName))
	gf.HeaderComment("version: " + compilerVersion)
	gf.HeaderComment(fmt.Sprintf("source: %s", file.Proto.GetName()))

	return &Generator{
		plugin:        plugin,
		file:          file,
		version:       compilerVersion,
		extensionName: extensionName,
		jen:           gf,
	}
}

// GetMessageOptionInfos retrieves the option information for a given proto extension.
//
// It takes a pointer to a protoimpl.ExtensionInfo as a parameter.
// The function iterates over the message types in the file and checks if the options contain the given extension.
// If the extension is found, the function retrieves the option and appends it to the result slice.
// The function returns the result slice containing the option information and a nil error.
func (g *Generator) GetMessageOptionInfos(e *protoimpl.ExtensionInfo) (map[string]proto.Message, error) {
	result := make(map[string]proto.Message)
	for _, fd := range g.file.Proto.GetMessageType() {
		if !proto.HasExtension(fd.GetOptions(), e) {
			continue
		}

		opt, ok := getOption(fd.Options.ProtoReflect(), e)
		if ok {
			result[fd.GetName()] = opt
		}
	}

	return result, nil
}

// GetEnumInfos returns the list of protogen.Enums from the Generator.
// It returns a slice of pointers to protogen.Enum.
func (g *Generator) GetEnumInfos() []*protogen.Enum {
	return g.file.Enums

}

func (g *Generator) GetMethodDescInfos() map[string][]*descriptorpb.MethodDescriptorProto {
	result := make(map[string][]*descriptorpb.MethodDescriptorProto)
	for _, svc := range g.file.Proto.GetService() {
		svcName := svc.GetName()
		for _, m := range svc.Method {
			if _, ok := result[svcName]; !ok {
				result[svc.GetName()] = make([]*descriptorpb.MethodDescriptorProto, 0)
			}
			result[svc.GetName()] = append(result[svcName], m)
		}
	}

	return result
}

// GetJen returns the Jen file from the Generator.
// Returns a pointer to a Jen file.
func (g *Generator) GetJen() *jen.File {
	return g.jen
}

// Render generates the Go code for the given Generator and writes it to a file.
//
// It constructs the filename using the Generator's file's generated filename prefix and the extension name.
// It creates a new generated file using the plugin's NewGeneratedFile method with the constructed filename and the file's Go import path.
// It renders the Generator's Jen file into the generated file using the Render method of the Jen file.
// If there is an error during rendering, it returns the error.
// Otherwise, it returns nil.
func (g *Generator) Render() error {
	filename := fmt.Sprintf("%s_%s.pb.go", g.file.GeneratedFilenamePrefix, g.extensionName)
	genFile := g.plugin.NewGeneratedFile(filename, g.file.GoImportPath)

	if err := g.jen.Render(genFile); err != nil {
		return err
	}

	return nil
}
