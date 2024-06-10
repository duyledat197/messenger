package main

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"

	"openmyth/messgener/protoc/generator"
)

// generateHTTPClientByService generates a client for a gRPC service.
// It takes a pointer to a protogen.Service and a pointer to a descriptorpb.ServiceDescriptorProto as parameters.
// It returns a pointer to a jen.Statement.
func generateHTTPClientByService(service *protogen.Service, protoService *descriptorpb.ServiceDescriptorProto) *jen.Statement {
	result := new(jen.Statement)

	serviceName := service.GoName

	result.Add(getClientStruct(serviceName))
	result.Add(getNewClient(serviceName))

	for i, m := range service.Methods {
		opt, err := generator.ExtractAPIOptions(protoService.Method[i])
		if opt == nil || err != nil {
			continue
		}

		methodName := m.GoName
		reqName := m.Input.GoIdent.GoName
		respName := m.Output.GoIdent.GoName

		httpMethod, path := generator.GetPathFromHTTPRule(opt)

		result.Add(getMethod(methodName, serviceName, reqName, respName, path, httpMethod))
	}

	return result
}

func getClientName(serviceName string) string {
	return fmt.Sprintf("%sHTTPClient", serviceName)
}

// getClientStruct generates a client struct for the given service.
func getClientStruct(serviceName string) *jen.Statement {
	return jen.Comment("HTTPClient is a http client for the "+serviceName+" service").Line().Type().Id(getClientName(serviceName)).Struct(
		jen.Id("BaseURL").String(),
		jen.Id("roundTripper").Qual("net/http", "RoundTripper"),
	).Line()
}

func getNewClient(serviceName string) *jen.Statement {
	return jen.Func().Id("New" + getClientName(serviceName)).
		Params(jen.Id("baseURL").String()).Params(jen.Op("*").Id(getClientName(serviceName))).BlockFunc(func(g *jen.Group) {
		g.Id("return &").Id(getClientName(serviceName)).Values(jen.Dict{
			jen.Id("BaseURL"):      jen.Id("baseURL"),
			jen.Id("roundTripper"): jen.Qual("openmyth/messgener/util/http_client", "NewRoundTripper").Call(),
		})
	})
}

// getClientFunc generates a client function for the given service.
func getClientFunc(serviceName string) *jen.Statement {
	return jen.Func().Params(jen.Id("c").Op("*").Id(serviceName)).Line()
}

const importPattern = `
	import (
		"context"
		"fmt"
		"net/url"

		"openmyth/messgener/util"

		"google.golang.org/grpc/status"
		"google.golang.org/grpc/codes"
	)
`

func getMethod(methodName, serviceName, reqName, respName, path, httpMethod string) *jen.Statement {
	return jen.Comment(methodName+" is a http call method for the "+serviceName+" service").Line().
		Func().Params(
		jen.Id("c").Op("*").Id(getClientName(serviceName)),
	).Id(methodName).
		Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("reqData").Op("*").Id(reqName)).
		Params(jen.Op("*").Id(respName), jen.Error()).BlockFunc(func(g *jen.Group) {

		g.List(jen.Id("path"), jen.Id("err")).Op(":=").Qual("net/url", "JoinPath").Call(jen.Id("c.BaseURL"), jen.Lit(path))

		g.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Id("return").List(jen.Nil(), jen.Qual("google.golang.org/grpc/status", "Errorf").
				Call(jen.Qual("google.golang.org/grpc/codes", "InvalidArgument"), jen.Qual("fmt", "Errorf").Call(jen.Lit("path is not valid: %w"), jen.Id("err")).Dot("Error").Call())),
		).Line()

		g.Id("reqClient, err").Op(":=").Qual("openmyth/messgener/util", "EncodeHTTPRequest").Call(jen.Id("ctx"), jen.Id("path"), jen.Lit(httpMethod), jen.Id("reqData"))
		g.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Id("return").List(jen.Nil(), jen.Qual("google.golang.org/grpc/status", "Errorf").
				Call(jen.Qual("google.golang.org/grpc/codes", "InvalidArgument"), jen.Qual("fmt", "Errorf").Call(jen.Lit("unable to encode http request: %w"), jen.Id("err")).Dot("Error").Call())),
		)
		g.Id("client").Op(":=").Qual("net/http", "Client").Values(jen.Dict{
			jen.Id("Transport"): jen.Id("c").Dot("roundTripper"),
		})

		g.Id("resp, err").Op(":=").Id("client").Dot("Do").Call(jen.Id("reqClient"))
		g.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Id("return").List(jen.Nil(), jen.Qual("google.golang.org/grpc/status", "Errorf").
				Call(jen.Qual("google.golang.org/grpc/codes", "Internal"), jen.Qual("fmt", "Errorf").Call(jen.Lit("unable to request: %w"), jen.Id("err")).Dot("Error").Call())),
		).Line()
		g.Id("return").List(jen.Qual("openmyth/messgener/util", "DecodeHTTPResponse").Types(jen.Id(respName)).Call(jen.Id("resp")))
	})

}

const methodPattern = `
// %s is a http call method for the %s service
func (c *%[2]sHTTPClient) %[1]s(ctx context.Context, reqData *%[3]s) (*%[4]s, error) {
	path, err := url.JoinPath(c.BaseURL, "%s")
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Errorf("path is not valid: %%w",err).Error())
	}
	reqClient, err := util.EncodeHTTPRequest(ctx,path,"%s",reqData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Errorf("unable to encode http request: %%w",err).Error())
	}

	client := http.Client{
		Transport: c.roundTripper,
	}

	resp, err := client.Do(reqClient)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Errorf("unable to request: %%w",err).Error())
	}

	return util.DecodeHTTPResponse[%[4]s](resp)
}
`
