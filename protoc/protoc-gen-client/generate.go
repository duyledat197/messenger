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
	result.Add(jen.Id(importPattern))
	result.Add(getClientStruct(serviceName))

	for i, m := range service.Methods {
		opt, err := generator.ExtractAPIOptions(protoService.Method[i])
		if opt == nil || err != nil {
			continue
		}

		funcName := m.GoName
		reqName := m.Input.GoIdent.GoName
		respName := m.Output.GoIdent.GoName

		httpMethod, path := generator.GetPathFromHTTPRule(opt)

		methodStr := fmt.Sprintf(methodPattern, funcName, serviceName, reqName, respName, path, httpMethod)

		result.Add(jen.Id(methodStr))
	}

	return result
}

// getClientStruct generates a client struct for the given service.
func getClientStruct(serviceName string) *jen.Statement {
	return jen.Comment("HTTPClient is a http client for the " + serviceName + " service").Line().Type().Id(fmt.Sprintf("%sHTTPClient", serviceName)).Struct(jen.Id("BaseURL").String())
}

// getClientFunc generates a client function for the given service.
func getClientFunc(serviceName string) *jen.Statement {
	return jen.Func().Params(jen.Id("c").Op("*").Id(serviceName))
}

const importPattern = `
	import (
		"context"
		"fmt"
		"net/url"
		"net/http"

		"openmyth/messgener/util"

		"google.golang.org/grpc/status"
		"google.golang.org/grpc/codes"
	)
`

const methodPattern = `
// %s is a http call method for the %s service
func (c *%[2]sHTTPClient) %[1]s(ctx context.Context, reqData *%[3]s) (*%[4]s, error) {
	path, err := url.JoinPath(c.BaseURL, "%s")
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Errorf("path is not valid: %%v",err.Error()).Error())
	}
	reqClient, err := util.EncodeHTTPRequest(ctx,path,"%s",reqData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Errorf("unable to encode http request: %%v",err.Error()).Error())
	}

	client := http.Client{}

	resp, err := client.Do(reqClient)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Errorf("unable to request: %%v",err.Error()).Error())
	}

	return util.DecodeHTTPResponse[%[4]s](resp)
}
`
