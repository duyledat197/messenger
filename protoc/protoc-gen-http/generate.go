package main

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/descriptorpb"
)

// generateConstAPIByService generates a constant statement for a service and its methods.
func generateConstAPIByService(serviceName string, methods []*descriptorpb.MethodDescriptorProto) *jen.Statement {
	var (
		stmts          []jen.Code
		methodKeyStmts []jen.Code
		result         *jen.Statement
	)

	for _, m := range methods {
		opts, err := extractAPIOptions(m)
		if err != nil {
			grpclog.Fatalf("unable to extract API options: %v", err)
		}
		if opts == nil {
			continue
		}
		stmts = append(stmts, generateConstAPIKey(serviceName, m.GetName(), opts))

		methodKeyStmts = append(methodKeyStmts, jen.Id(getAPIKey(serviceName, m.GetName())))
	}

	if len(stmts) > 0 {
		result = jen.Const().Defs(stmts...).Line()
		result.Var().Id(fmt.Sprintf("%s_APIs", serviceName)).Op("=").Index().String().BlockFunc(
			func(g *jen.Group) {
				for _, stmt := range methodKeyStmts {
					g.Add(stmt).Op(",")
				}
			},
		)
	}

	return result
}

// getAPIKey generates an API key by concatenating the service name and method name with an underscore.
func getAPIKey(serviceName, methodName string) string {
	return fmt.Sprintf("%s_%s_API", serviceName, methodName)
}

// generateConstAPIKey generates a constant statement for an API key by concatenating the service name and method name with an underscore.
func generateConstAPIKey(serviceName, methodName string, opts *options.HttpRule) *jen.Statement {
	apiName := getAPIKey(serviceName, methodName)
	return jen.Id(apiName).Op("=").Lit(getPathFromHTTPRule(opts))
}
