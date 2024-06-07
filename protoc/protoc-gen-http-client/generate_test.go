package main

import (
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/require"
)

func Test_getClientStruct(t *testing.T) {
	gf := jen.NewFile("test")
	got := getClientStruct("Hello")

	expected := `type HelloHTTPClient struct {
	BaseURL string
}`

	buf := bytes.Buffer{}
	err := gf.Add(got).Render(&buf)
	require.NoError(t, err)
	require.Equal(t, expected, buf.String())
}
