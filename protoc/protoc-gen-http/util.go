package main

import (
	"fmt"
	"regexp"
	"strings"

	options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var re = regexp.MustCompile(`\{(.*?)\}`)

// getPathFromHTTPRule returns the HTTP method and path extracted from the given HttpRule.
// The function takes a pointer to an HttpRule as input and returns a string representing the HTTP method and path.
// The HttpRule is a protocol buffer message that defines the mapping between an HTTP REST method and an API method.
// The function uses a type switch to determine the type of the HttpRule pattern and extracts the corresponding HTTP method and path.
// If the pattern is a Get, Post, Put, or Delete pattern, the function returns the corresponding HTTP method and path with wildcards replaced by the actual values.
// If the pattern is none of the above, the function returns an empty string.

func getPathFromHTTPRule(opt *options.HttpRule) string {
	switch opt.Pattern.(type) {
	case *options.HttpRule_Get:
		rule := opt.Pattern.(*options.HttpRule_Get)
		return "GET " + re.ReplaceAllString(rule.Get, "*")
	case *options.HttpRule_Post:
		rule := opt.Pattern.(*options.HttpRule_Post)
		return "POST " + re.ReplaceAllString(rule.Post, "*")
	case *options.HttpRule_Put:
		rule := opt.Pattern.(*options.HttpRule_Put)
		return "PUT " + re.ReplaceAllString(rule.Put, "*")
	case *options.HttpRule_Delete:
		rule := opt.Pattern.(*options.HttpRule_Delete)
		return "DELETE " + re.ReplaceAllString(rule.Delete, "*")
	}

	return ""
}

// normalizeFullname description of the Go function.
func normalizeFullname(fn protoreflect.FullName) string {
	return strings.ReplaceAll(string(fn), ".", "_")
}

// extractAPIOptions retrieves the HTTP options for a given method descriptor.
func extractAPIOptions(meth *descriptorpb.MethodDescriptorProto) (*options.HttpRule, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, options.E_Http) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, options.E_Http)
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}

	return opts, nil
}
