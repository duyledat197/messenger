// Package errcode ...
package errcode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetrieveTranslate(t *testing.T) {
	require.Equal(t, "not found", RetrieveTranslate("not_found", "en"))
	require.Equal(t, "không tìm thấy", RetrieveTranslate("not_found", "vi"))
}
