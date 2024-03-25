package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type foo struct {
	Name     string
	Age      int
	UserInfo any
}

func (aa *foo) TableName() string {
	return "foo"
}

func TestFieldMap(t *testing.T) {

	got, _ := FieldMap(&foo{})

	want := []string{"name", "age", "user_info"}

	require.Equal(t, want, got)
}
